// 实现功能： 从nsq里取得日志，汇总后插入到数据库中
package consumer

import (
	"log"
	"net"
	"os"
	"tencent/logagent/prod"
	"time"

	"github.com/bitly/nsq/util"
	"github.com/bitly/nsq/util/lookupd"
)

const (
	INIT_READER_NUM = 80 // 日志reader的个数认为不会超过这个值
)

type ConsServer struct {
	waitGroup util.WaitGroupWrapper
	opts      *ConsOption
	quit      chan bool

	jsonconf  *prod.JsonConf
	jsonMtime time.Time

	readers []*logReader

	dbpool *DBPool

	httpListener net.Listener
}

func NewConsServer(opt *ConsOption) *ConsServer {
	p := new(ConsServer)
	p.opts = opt
	p.quit = make(chan bool)

	p.readers = make([]*logReader, 0, INIT_READER_NUM)

	p.dbpool = NewDBPool(opt.SqlDir)

	return p
}

func (s *ConsServer) Main() {
	httpListener, err := net.Listen("tcp", s.opts.HttpAddr)
	if err != nil {
		log.Fatalf("ERR: listen (%s) failed - %s", s.opts.HttpAddr, err.Error())
	}
	s.httpListener = httpListener
	httpServer := &httpServer{context: &Context{s}}
	s.waitGroup.Wrap(func() { util.HTTPServer(s.httpListener, httpServer, "HTTP") })

	cb := func() {
		s.loadConf()

		ticker := time.NewTicker(time.Minute * 15)
		ticker2 := time.NewTicker(time.Minute)

		for {
			select {
			case <-ticker.C:
				s.reloadJson()

			case <-ticker2.C:
				s.updateReader(s.dbpool.InputChan)

			case <-s.dbpool.FailChan:
				s.Stop()
				log.Fatal("Fatal: pool notwork, server exit")

			case <-s.quit:
				goto END
			}
		}

	END:
		log.Println("INF: server quit loop")
		s.Stop()

		if s.httpListener != nil {
			s.httpListener.Close()
		}
	}

	log.Println("INF: server run")
	s.waitGroup.Wrap(cb)

}

func (s *ConsServer) Exit() {
	s.quit <- true
	s.waitGroup.Wait()
	log.Println("INF: server exit")
}

func (s *ConsServer) loadConf() {
	log.Printf("INF, loadConf")

	s.jsonconf = prod.NewJsonConf(s.opts.JsFile)
	fi, _ := os.Stat(s.opts.JsFile)
	s.jsonMtime = fi.ModTime()

	s.dbpool.Start(s.jsonconf)
	s.updateReader(s.dbpool.InputChan)
}

func (s *ConsServer) updateReader(sqlChan chan SqlItem) {
	topics, err := lookupd.GetLookupdTopics(s.opts.NsqlookupdHttpAddr)
	if err != nil {
		log.Fatal("ERR, lookupdlGetLookupdTopics, %s, %s", s.opts.NsqlookupdHttpAddr, err)
	}

	for i := 0; i < len(s.jsonconf.Log_list); i++ {
		table := s.jsonconf.Log_list[i].Table

		found := false
		for j := 0; j < len(topics); j++ {
			if table == topics[j] {
				found = true
				break
			}
		}
		if !found {
			log.Printf("WAR: topic %s not in NSQD", table)
			continue
		}

		idx := -1
		for k := 0; k < len(s.readers); k++ {
			if s.readers[k].table == table {
				idx = k
				break
			}
		}
		if idx != -1 {
			if !s.readers[idx].IsRun() {
				log.Printf("INF, restart reader[%s]", table)
				s.readers[idx].Start(s.jsonconf, s.opts.MaxInFlight, s.opts.NsqlookupdHttpAddr, sqlChan)
			}
			continue
		}

		r := NewlogReader(table, s.opts.Duration)
		s.readers = append(s.readers, r)
		r.Start(s.jsonconf, s.opts.MaxInFlight, s.opts.NsqlookupdHttpAddr, sqlChan)
		log.Printf("INF: new reader[%s]", table)
	}
}

func (s *ConsServer) reloadJson() {
	fi, _ := os.Stat(s.opts.JsFile)
	if fi.ModTime() == s.jsonMtime {
		// 配置文件没有修改
		return
	}

	log.Println("INF, reloadJson")
	s.Stop()

	s.loadConf()
}

func (s *ConsServer) Stop() {
	// 先停止所有的输入，然后再吐出数据
	for _, r := range s.readers {
		r.Stop()
	}
	for _, r := range s.readers {
		r.Clear()
	}

	s.dbpool.Stop()
}
