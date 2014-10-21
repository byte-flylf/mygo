// 测试：基于nsq的实现
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"os"
	"time"
)

type response struct {
	Time     int    `json:"time"`
	Hostname string `json:'hostname'`
	Sum      int    `json: sum`
	Status   string `json: status`
	Message  string `json: message`
}

var handlers = map[string](func(*simplejson.Json) response){
	"get_time":     GetTime,
	"get_hostname": GetHostName,
	"get_sum":      GetSum,
}

func GetTime(js *simplejson.Json) response {
	return response{Time: int(time.Now().Nanosecond())}
}

func GetHostName(js *simplejson.Json) response {
	return response{Hostname: "samhost"}
}

func GetSum(js *simplejson.Json) response {
	a, _ := js.Get("a").Int()
	b, _ := js.Get("b").Int()
	return response{Sum: a + b}
}

func SendMessage(topic string, body []byte) {
	httpclient := &http.Client{}
	endpoint := fmt.Sprintf("http://127.0.0.1:4151/put?topic=%s", topic)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	resp, err := httpclient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp.Body.Close()

}

type MyHandler struct {
	q               *nsq.Consumer
	messageSent     int
	messageReceived int
	messagesFailed  int
}

func (h *MyHandler) LogFailedMessage(message *nsq.Message) {
	h.messagesFailed++
	h.q.Stop()
}

func (h *MyHandler) HandleMessage(message *nsq.Message) error {
	fmt.Println("RPC <-- ", string(message.Body))

	request, err := simplejson.NewJson(message.Body)
	if err != nil {
		return err
	}
	call, err := request.Get("call").String()
	var res response
	if err != nil {
		res.Status = "error"
		res.Message = "You need to specify a call"
	} else {
		if f, ok := handlers[call]; !ok {
			res.Status = "error"
			res.Message = "call not found"
		} else {
			res = f(request)
		}
	}
	fmt.Println(&res)
	b, err := json.Marshal(&res)
	if err != nil {
		fmt.Println("json err", err)
		return err
	}
	fmt.Println("marshal succ", string(b), res)
	from, err := request.Get("from").String()
	if err != nil {
		fmt.Println("param from", err)
	}
	SendMessage(from, b)

	return nil
}

func main() {
	cfg := nsq.NewConfig()
	cfg.Set("verbose", true)
	reader, err := nsq.NewConsumer("rpc", "_", cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	reader.SetLogger(log.New(os.Stderr, "", log.LstdFlags), nsq.LogLevelInfo)

	h := &MyHandler{q: reader}
	reader.SetHandler(h)
	addr := "127.0.0.1:4150"
	err = reader.ConnectToNSQD(addr)
	if err != nil {
		log.Fatalf(err.Error())
	}
	<-reader.StopChan

}
