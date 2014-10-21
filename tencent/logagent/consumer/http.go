package consumer

import (
	"errors"
	"fmt"
	"github.com/bitly/nsq/util"
	"log"
	"net/http"
	httpprof "net/http/pprof"
)

type httpServer struct {
	context *Context
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := s.debugRouter(w, req)
	if err != nil {
		log.Printf("ERROR: %s", err)
		util.ApiResponse(w, 404, "NOT_FOUND", nil)
	}
}

func (s *httpServer) debugRouter(w http.ResponseWriter, req *http.Request) error {
	switch req.URL.Path {
	case "/debug/pprof":
		httpprof.Index(w, req)
	case "/debug/pprof/cmdline":
		httpprof.Cmdline(w, req)
	case "/debug/pprof/symbol":
		httpprof.Symbol(w, req)
	case "/debug/pprof/heap":
		httpprof.Handler("heap").ServeHTTP(w, req)
	case "/debug/pprof/goroutine":
		httpprof.Handler("goroutine").ServeHTTP(w, req)
	case "/debug/pprof/profile":
		httpprof.Profile(w, req)
	case "/debug/pprof/block":
		httpprof.Handler("block").ServeHTTP(w, req)
	case "/debug/pprof/threadcreate":
		httpprof.Handler("threadcreate").ServeHTTP(w, req)
	default:
		return errors.New(fmt.Sprintf("404 %s", req.URL.Path))
	}
	return nil
}
