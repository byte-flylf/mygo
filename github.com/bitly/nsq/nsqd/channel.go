package nsqd

import (
	"bytes"
	"container/heap"
	"errors"
	"log"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bitly/nsq/util"
	"github.com/bitly/nsq/util/pqueue"
)

// the amount of time a worker will wait when idle
const defaultWorkerWait = 100 * time.Millisecond

type Consumer interface {
	UnPause()
	Pause()
	Close() error
	TimedOutMessage()
	Stats() ClientStats
	Empty()
}

// Channel represents the concrete type for a NSQ channel (and also
// implements the Queue interface)
//
// There can be multiple channels per topic, each with there own unique set
// of subscribers (clients).
//
// Channels maintain all client and message metadata, orchestrating in-flight
// messages, timeouts, requeuing, etc.
type Channel struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	requeueCount uint64
	messageCount uint64
	timeoutCount uint64

	sync.RWMutex

	topicName string
	name      string
	context   *context

	backend BackendQueue

	incomingMsgChan chan *Message
	memoryMsgChan   chan *Message
	clientMsgChan   chan *Message
	exitChan        chan int
	waitGroup       util.WaitGroupWrapper
	exitFlag        int32

	// state tracking
	clients          map[int64]Consumer
	paused           int32
	ephemeralChannel bool
	deleteCallback   func(*Channel)
	deleter          sync.Once

	// Stats tracking
	e2eProcessingLatencyStream *util.Quantile

	// TODO: these can be DRYd up
	deferredMessages map[MessageID]*pqueue.Item
	deferredPQ       pqueue.PriorityQueue
	deferredMutex    sync.Mutex
	inFlightMessages map[MessageID]*Message
	inFlightPQ       inFlightPqueue
	inFlightMutex    sync.Mutex

	// stat counters
	bufferedCount int32
}

// NewChannel creates a new instance of the Channel type and returns a pointer
func NewChannel(topicName string, channelName string, context *context,
	deleteCallback func(*Channel)) *Channel {

	c := &Channel{
		topicName:       topicName,
		name:            channelName,
		incomingMsgChan: make(chan *Message, 1),
		memoryMsgChan:   make(chan *Message, context.nsqd.options.MemQueueSize),
		clientMsgChan:   make(chan *Message),
		exitChan:        make(chan int),
		clients:         make(map[int64]Consumer),
		deleteCallback:  deleteCallback,
		context:         context,
	}
	if len(context.nsqd.options.E2EProcessingLatencyPercentiles) > 0 {
		c.e2eProcessingLatencyStream = util.NewQuantile(
			context.nsqd.options.E2EProcessingLatencyWindowTime,
			context.nsqd.options.E2EProcessingLatencyPercentiles,
		)
	}

	c.initPQ()

	if strings.HasSuffix(channelName, "#ephemeral") {
		c.ephemeralChannel = true
		c.backend = newDummyBackendQueue()
	} else {
		// backend names, for uniqueness, automatically include the topic...
		backendName := getBackendName(topicName, channelName)
		c.backend = newDiskQueue(backendName,
			context.nsqd.options.DataPath,
			context.nsqd.options.MaxBytesPerFile,
			context.nsqd.options.SyncEvery,
			context.nsqd.options.SyncTimeout)
	}

	go c.messagePump()

	c.waitGroup.Wrap(func() { c.router() })
	c.waitGroup.Wrap(func() { c.deferredWorker() })
	c.waitGroup.Wrap(func() { c.inFlightWorker() })

	go c.context.nsqd.Notify(c)

	return c
}

func (c *Channel) initPQ() {
	pqSize := int(math.Max(1, float64(c.context.nsqd.options.MemQueueSize)/10))

	c.inFlightMessages = make(map[MessageID]*Message)
	c.deferredMessages = make(map[MessageID]*pqueue.Item)

	c.inFlightMutex.Lock()
	c.inFlightPQ = newInFlightPqueue(pqSize)
	c.inFlightMutex.Unlock()

	c.deferredMutex.Lock()
	c.deferredPQ = pqueue.New(pqSize)
	c.deferredMutex.Unlock()
}

// Exiting returns a boolean indicating if this channel is closed/exiting
func (c *Channel) Exiting() bool {
	return atomic.LoadInt32(&c.exitFlag) == 1
}

// Delete empties the channel and closes
func (c *Channel) Delete() error {
	return c.exit(true)
}

// Close cleanly closes the Channel
func (c *Channel) Close() error {
	return c.exit(false)
}

func (c *Channel) exit(deleted bool) error {
	if !atomic.CompareAndSwapInt32(&c.exitFlag, 0, 1) {
		return errors.New("exiting")
	}

	if deleted {
		log.Printf("CHANNEL(%s): deleting", c.name)

		// since we are explicitly deleting a channel (not just at system exit time)
		// de-register this from the lookupd
		go c.context.nsqd.Notify(c)
	} else {
		log.Printf("CHANNEL(%s): closing", c.name)
	}

	// this forceably closes client connections
	c.RLock()
	for _, client := range c.clients {
		client.Close()
	}
	c.RUnlock()

	close(c.exitChan)

	// handle race condition w/ things writing into incomingMsgChan
	c.Lock()
	close(c.incomingMsgChan)
	c.Unlock()

	// synchronize the close of router() and pqWorkers (2)
	c.waitGroup.Wait()

	if deleted {
		// empty the queue (deletes the backend files, too)
		c.Empty()
		return c.backend.Delete()
	}

	// write anything leftover to disk
	c.flush()
	return c.backend.Close()
}

func (c *Channel) Empty() error {
	c.Lock()
	defer c.Unlock()

	c.initPQ()
	for _, client := range c.clients {
		client.Empty()
	}

	clientMsgChan := c.clientMsgChan
	for {
		select {
		case _, ok := <-clientMsgChan:
			if !ok {
				// c.clientMsgChan may be closed while in this loop
				// so just remove it from the select so we can make progress
				clientMsgChan = nil
			}
		case <-c.memoryMsgChan:
		default:
			goto finish
		}
	}

finish:
	return c.backend.Empty()
}

// flush persists all the messages in internal memory buffers to the backend
// it does not drain inflight/deferred because it is only called in Close()
func (c *Channel) flush() error {
	var msgBuf bytes.Buffer

	// messagePump is responsible for closing the channel it writes to
	// this will read until its closed (exited)
	for msg := range c.clientMsgChan {
		log.Printf("CHANNEL(%s): recovered buffered message from clientMsgChan", c.name)
		writeMessageToBackend(&msgBuf, msg, c.backend)
	}

	if len(c.memoryMsgChan) > 0 || len(c.inFlightMessages) > 0 || len(c.deferredMessages) > 0 {
		log.Printf("CHANNEL(%s): flushing %d memory %d in-flight %d deferred messages to backend",
			c.name, len(c.memoryMsgChan), len(c.inFlightMessages), len(c.deferredMessages))
	}

	for {
		select {
		case msg := <-c.memoryMsgChan:
			err := writeMessageToBackend(&msgBuf, msg, c.backend)
			if err != nil {
				log.Printf("ERROR: failed to write message to backend - %s", err.Error())
			}
		default:
			goto finish
		}
	}

finish:
	for _, msg := range c.inFlightMessages {
		err := writeMessageToBackend(&msgBuf, msg, c.backend)
		if err != nil {
			log.Printf("ERROR: failed to write message to backend - %s", err.Error())
		}
	}

	for _, item := range c.deferredMessages {
		msg := item.Value.(*Message)
		err := writeMessageToBackend(&msgBuf, msg, c.backend)
		if err != nil {
			log.Printf("ERROR: failed to write message to backend - %s", err.Error())
		}
	}

	return nil
}

func (c *Channel) Depth() int64 {
	return int64(len(c.memoryMsgChan)) + c.backend.Depth() + int64(atomic.LoadInt32(&c.bufferedCount))
}

func (c *Channel) Pause() error {
	return c.doPause(true)
}

func (c *Channel) UnPause() error {
	return c.doPause(false)
}

func (c *Channel) doPause(pause bool) error {
	if pause {
		atomic.StoreInt32(&c.paused, 1)
	} else {
		atomic.StoreInt32(&c.paused, 0)
	}

	c.RLock()
	for _, client := range c.clients {
		if pause {
			client.Pause()
		} else {
			client.UnPause()
		}
	}
	c.RUnlock()

	c.context.nsqd.Lock()
	defer c.context.nsqd.Unlock()
	// pro-actively persist metadata so in case of process failure
	// nsqd won't suddenly (un)pause a channel
	return c.context.nsqd.PersistMetadata()
}

func (c *Channel) IsPaused() bool {
	return atomic.LoadInt32(&c.paused) == 1
}

// PutMessage writes to the appropriate incoming message channel
// (which will be routed asynchronously)
func (c *Channel) PutMessage(msg *Message) error {
	c.RLock()
	defer c.RUnlock()
	if atomic.LoadInt32(&c.exitFlag) == 1 {
		return errors.New("exiting")
	}
	c.incomingMsgChan <- msg
	atomic.AddUint64(&c.messageCount, 1)
	return nil
}

// TouchMessage resets the timeout for an in-flight message
func (c *Channel) TouchMessage(clientID int64, id MessageID, clientMsgTimeout time.Duration) error {
	msg, err := c.popInFlightMessage(clientID, id)
	if err != nil {
		return err
	}
	c.removeFromInFlightPQ(msg)

	newTimeout := time.Now().Add(clientMsgTimeout)
	if newTimeout.Sub(msg.deliveryTS) >=
		c.context.nsqd.options.MaxMsgTimeout {
		// we would have gone over, set to the max
		newTimeout = msg.deliveryTS.Add(c.context.nsqd.options.MaxMsgTimeout)
	}

	msg.pri = newTimeout.UnixNano()
	err = c.pushInFlightMessage(msg)
	if err != nil {
		return err
	}
	c.addToInFlightPQ(msg)
	return nil
}

// FinishMessage successfully discards an in-flight message
func (c *Channel) FinishMessage(clientID int64, id MessageID) error {
	msg, err := c.popInFlightMessage(clientID, id)
	if err != nil {
		return err
	}
	c.removeFromInFlightPQ(msg)
	if c.e2eProcessingLatencyStream != nil {
		c.e2eProcessingLatencyStream.Insert(msg.Timestamp)
	}
	return nil
}

// RequeueMessage requeues a message based on `time.Duration`, ie:
//
// `timeoutMs` == 0 - requeue a message immediately
// `timeoutMs`  > 0 - asynchronously wait for the specified timeout
//     and requeue a message (aka "deferred requeue")
//
func (c *Channel) RequeueMessage(clientID int64, id MessageID, timeout time.Duration) error {
	// remove from inflight first
	msg, err := c.popInFlightMessage(clientID, id)
	if err != nil {
		return err
	}
	c.removeFromInFlightPQ(msg)

	if timeout == 0 {
		return c.doRequeue(msg)
	}

	// deferred requeue
	return c.StartDeferredTimeout(msg, timeout)
}

// AddClient adds a client to the Channel's client list
func (c *Channel) AddClient(clientID int64, client Consumer) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.clients[clientID]
	if ok {
		return
	}
	c.clients[clientID] = client
}

// RemoveClient removes a client from the Channel's client list
func (c *Channel) RemoveClient(clientID int64) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.clients[clientID]
	if !ok {
		return
	}
	delete(c.clients, clientID)

	if len(c.clients) == 0 && c.ephemeralChannel == true {
		go c.deleter.Do(func() { c.deleteCallback(c) })
	}
}

func (c *Channel) StartInFlightTimeout(msg *Message, clientID int64, timeout time.Duration) error {
	now := time.Now()
	msg.clientID = clientID
	msg.deliveryTS = now
	msg.pri = now.Add(timeout).UnixNano()
	err := c.pushInFlightMessage(msg)
	if err != nil {
		return err
	}
	c.addToInFlightPQ(msg)
	return nil
}

func (c *Channel) StartDeferredTimeout(msg *Message, timeout time.Duration) error {
	absTs := time.Now().Add(timeout).UnixNano()
	item := &pqueue.Item{Value: msg, Priority: absTs}
	err := c.pushDeferredMessage(item)
	if err != nil {
		return err
	}
	c.addToDeferredPQ(item)
	return nil
}

// doRequeue performs the low level operations to requeue a message
func (c *Channel) doRequeue(msg *Message) error {
	c.RLock()
	defer c.RUnlock()
	if atomic.LoadInt32(&c.exitFlag) == 1 {
		return errors.New("exiting")
	}
	c.incomingMsgChan <- msg
	atomic.AddUint64(&c.requeueCount, 1)
	return nil
}

// pushInFlightMessage atomically adds a message to the in-flight dictionary
func (c *Channel) pushInFlightMessage(msg *Message) error {
	c.Lock()
	_, ok := c.inFlightMessages[msg.ID]
	if ok {
		c.Unlock()
		return errors.New("ID already in flight")
	}
	c.inFlightMessages[msg.ID] = msg
	c.Unlock()
	return nil
}

// popInFlightMessage atomically removes a message from the in-flight dictionary
func (c *Channel) popInFlightMessage(clientID int64, id MessageID) (*Message, error) {
	c.Lock()
	msg, ok := c.inFlightMessages[id]
	if !ok {
		c.Unlock()
		return nil, errors.New("ID not in flight")
	}
	if msg.clientID != clientID {
		c.Unlock()
		return nil, errors.New("client does not own message")
	}
	delete(c.inFlightMessages, id)
	c.Unlock()
	return msg, nil
}

func (c *Channel) addToInFlightPQ(msg *Message) {
	c.inFlightMutex.Lock()
	c.inFlightPQ.Push(msg)
	c.inFlightMutex.Unlock()
}

func (c *Channel) removeFromInFlightPQ(msg *Message) {
	c.inFlightMutex.Lock()
	if msg.index == -1 {
		// this item has already been popped off the pqueue
		c.inFlightMutex.Unlock()
		return
	}
	c.inFlightPQ.Remove(msg.index)
	c.inFlightMutex.Unlock()
}

func (c *Channel) pushDeferredMessage(item *pqueue.Item) error {
	c.Lock()
	defer c.Unlock()

	// TODO: these map lookups are costly
	id := item.Value.(*Message).ID
	_, ok := c.deferredMessages[id]
	if ok {
		return errors.New("ID already deferred")
	}
	c.deferredMessages[id] = item

	return nil
}

func (c *Channel) popDeferredMessage(id MessageID) (*pqueue.Item, error) {
	c.Lock()
	defer c.Unlock()

	// TODO: these map lookups are costly
	item, ok := c.deferredMessages[id]
	if !ok {
		return nil, errors.New("ID not deferred")
	}
	delete(c.deferredMessages, id)

	return item, nil
}

func (c *Channel) addToDeferredPQ(item *pqueue.Item) {
	c.deferredMutex.Lock()
	defer c.deferredMutex.Unlock()

	heap.Push(&c.deferredPQ, item)
}

// Router handles the muxing of incoming Channel messages, either writing
// to the in-memory channel or to the backend
func (c *Channel) router() {
	var msgBuf bytes.Buffer
	for msg := range c.incomingMsgChan {
		select {
		case c.memoryMsgChan <- msg:
		default:
			err := writeMessageToBackend(&msgBuf, msg, c.backend)
			if err != nil {
				log.Printf("CHANNEL(%s) ERROR: failed to write message to backend - %s", c.name, err.Error())
				// theres not really much we can do at this point, you're certainly
				// going to lose messages...
			}
		}
	}

	log.Printf("CHANNEL(%s): closing ... router", c.name)
}

// messagePump reads messages from either memory or backend and writes
// to the client output go channel
//
// it is also performs in-flight accounting and initiates the auto-requeue
// goroutine
func (c *Channel) messagePump() {
	var msg *Message
	var buf []byte
	var err error

	for {
		// do an extra check for closed exit before we select on all the memory/backend/exitChan
		// this solves the case where we are closed and something else is draining clientMsgChan into
		// backend. we don't want to reverse that
		if atomic.LoadInt32(&c.exitFlag) == 1 {
			goto exit
		}

		select {
		case msg = <-c.memoryMsgChan:
		case buf = <-c.backend.ReadChan():
			msg, err = decodeMessage(buf)
			if err != nil {
				log.Printf("ERROR: failed to decode message - %s", err.Error())
				continue
			}
		case <-c.exitChan:
			goto exit
		}

		msg.Attempts++

		atomic.StoreInt32(&c.bufferedCount, 1)
		c.clientMsgChan <- msg
		atomic.StoreInt32(&c.bufferedCount, 0)
		// the client will call back to mark as in-flight w/ it's info
	}

exit:
	log.Printf("CHANNEL(%s): closing ... messagePump", c.name)
	close(c.clientMsgChan)
}

func (c *Channel) deferredWorker() {
	c.pqWorker(&c.deferredPQ, &c.deferredMutex, func(item *pqueue.Item) {
		msg := item.Value.(*Message)
		_, err := c.popDeferredMessage(msg.ID)
		if err != nil {
			return
		}
		c.doRequeue(msg)
	})
}

func (c *Channel) inFlightWorker() {
	ticker := time.NewTicker(defaultWorkerWait)
	for {
		select {
		case <-ticker.C:
		case <-c.exitChan:
			goto exit
		}
		now := time.Now().UnixNano()
		for {
			c.inFlightMutex.Lock()
			msg, _ := c.inFlightPQ.PeekAndShift(now)
			c.inFlightMutex.Unlock()

			if msg == nil {
				break
			}

			_, err := c.popInFlightMessage(msg.clientID, msg.ID)
			if err != nil {
				break
			}
			atomic.AddUint64(&c.timeoutCount, 1)
			c.RLock()
			client, ok := c.clients[msg.clientID]
			c.RUnlock()
			if ok {
				client.TimedOutMessage()
			}
			c.doRequeue(msg)
		}
	}

exit:
	log.Printf("CHANNEL(%s): closing ... inFlightWorker", c.name)
	ticker.Stop()
}

// generic loop (executed in a goroutine) that periodically wakes up to walk
// the priority queue and call the callback
func (c *Channel) pqWorker(pq *pqueue.PriorityQueue, mutex *sync.Mutex, callback func(item *pqueue.Item)) {
	ticker := time.NewTicker(defaultWorkerWait)
	for {
		select {
		case <-ticker.C:
		case <-c.exitChan:
			goto exit
		}
		now := time.Now().UnixNano()
		for {
			mutex.Lock()
			item, _ := pq.PeekAndShift(now)
			mutex.Unlock()

			if item == nil {
				break
			}

			callback(item)
		}
	}

exit:
	log.Printf("CHANNEL(%s): closing ... pqueue worker", c.name)
	ticker.Stop()
}
