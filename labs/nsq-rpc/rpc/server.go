// rpc on NSQ
package rpc

// todo: 如何实现广播?

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"log"
	"os"
)

type RpcServer struct {
	rd *nsq.Consumer
	wr *nsq.Producer
}

// 服务端:
// reader的topic应该是svr的名称, 服务启动就应该是确定的
// writer的topic由客户端来告知
func NewServer(svrname string, addr string) (*RpcServer, error) {
	svr := new(RpcServer)
	var err error

	cfg := nsq.NewConfig()
	cfg.Set("verbose", true)

	svr.rd, err = nsq.NewConsumer(svrname, "_", cfg)
	if err != nil {
		return nil, fmt.Errorf("fail to new consumer: %s", err)
	}
	svr.rd.SetLogger(log.New(os.Stderr, "", log.LstdFlags), nsq.LogLevelInfo)
	err = svr.rd.ConnectToNSQD(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to connected nsqlookupd", err)
	}

	svr.wr = nsq.NewProducer(addr, cfg)

	return svr, nil
}

func (s *RpcServer) Run() {
	<-s.rd.StopChan
}

type ServerHandler struct {
	q               *nsq.Consumer
	messageSent     int
	messageReceived int
	messagesFailed  int
}

func (h *ServerHandler) LogFailedMessage(message *nsq.Message) {
	h.messagesFailed++
	h.q.Stop()
}

func (h *ServerHandler) HandleMessage(message *nsq.Message) error {
	h.messageReceived++

	id, enc, topic, addr, body, err := decodeRequest(message.Body)
	if err != nil {
		h.messagesFailed++
		log.Printf("fail to decode request: %q", err)
		return err
	}
	log.Println("succ to decode", id, enc, topic, addr, body)
	if enc != h.q.enc {
		return fmt.Errorf("not supported protocol")
	}

	return
}
