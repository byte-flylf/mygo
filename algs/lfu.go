package algs

import (
	"container/list"
	"sync"
)

type LFUCache struct {
	// If len > UpperBound, cache will automatically evict
	// down to LowerBound.  If either value is 0, this behavior
	// is disabled.
	UpperBound int
	LowerBound int
	values     map[string]*cacheEntry
	freqs      *list.List
	len        int
	lock       *sync.Mutex
}

type cacheEntry struct {
	key      string
	value    interface{}
	freqNode *list.Element
}

type listEntry struct {
	entries map[*cacheEntry]struct{} // set
	freq    int
}

func NewLFUCache() *LFUCache {
	c := new(LFUCache)
	c.values = make(map[string]*cacheEntry)
	c.freqs = list.New()
	c.lock = new(sync.Mutex)
	return c
}

func (c *LFUCache) Get(key string) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value
	}
	return nil
}

func (c *LFUCache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		// value already exists for key.  overwrite
		e.value = value
		c.increment(e)
	} else {
		// value doesn't exist.  insert
		e := new(cacheEntry)
		e.key = key
		e.value = value
		c.values[key] = e
		c.increment(e)
		c.len++
		// bounds mgmt
		if c.UpperBound > 0 && c.LowerBound > 0 {
			if c.len > c.UpperBound {
				c.evit(c.len - c.LowerBound)
			}
		}
	}
}

func (c *LFUCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.len
}

func (c *LFUCache) Evict(count int) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.evit(count)
}

func (c *LFUCache) evit(count int) int {
	// No lock here so it can be called
	// from within the lock (during Set)
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry, _ := range place.Value.(*listEntry).entries {
				if i < count {
					delete(c.values, entry.key)
					c.remEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *LFUCache) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element

	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next()
	}
	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		// create a new list entry
		li := new(listEntry)
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry]struct{})
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry).entries[e] = struct{}{}
	if currentPlace != nil {
		// remove from current position
		c.remEntry(currentPlace, e)
	}
}

func (c *LFUCache) remEntry(place *list.Element, entry *cacheEntry) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}
