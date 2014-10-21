package nsqd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/bmizerany/assert"
)

func TestHTTPput(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_put" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	buf := bytes.NewBuffer([]byte("test message"))
	url := fmt.Sprintf("http://%s/put?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(1))
}

func TestHTTPputEmpty(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_put_empty" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	buf := bytes.NewBuffer([]byte(""))
	url := fmt.Sprintf("http://%s/put?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, resp.StatusCode, 500)
	assert.Equal(t, string(body), `{"status_code":500,"status_txt":"MSG_EMPTY","data":null}`)

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(0))
}

func TestHTTPmput(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_mput" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	msg := []byte("test message")
	msgs := make([][]byte, 0)
	for i := 0; i < 4; i++ {
		msgs = append(msgs, msg)
	}
	buf := bytes.NewBuffer(bytes.Join(msgs, []byte("\n")))

	url := fmt.Sprintf("http://%s/mput?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(4))
}

func TestHTTPmputEmpty(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_mput_empty" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	msg := []byte("test message")
	msgs := make([][]byte, 0)
	for i := 0; i < 4; i++ {
		msgs = append(msgs, msg)
	}
	buf := bytes.NewBuffer(bytes.Join(msgs, []byte("\n")))
	_, err := buf.Write([]byte("\n"))
	assert.Equal(t, err, nil)

	url := fmt.Sprintf("http://%s/mput?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(4))
}

func TestHTTPmputBinary(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_mput_bin" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	mpub := make([][]byte, 0)
	for i := 0; i < 5; i++ {
		mpub = append(mpub, make([]byte, 100))
	}
	cmd, _ := nsq.MultiPublish(topicName, mpub)
	buf := bytes.NewBuffer(cmd.Body)

	url := fmt.Sprintf("http://%s/mput?topic=%s&binary=true", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(5))
}

func TestHTTPSRequire(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	options := NewNSQDOptions()
	options.Verbose = true
	options.TLSCert = "./test/certs/server.pem"
	options.TLSKey = "./test/certs/server.key"
	options.TLSClientAuthPolicy = "require"
	_, httpAddr, nsqd := mustStartNSQD(options)

	defer nsqd.Exit()

	topicName := "test_http_put_req" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	buf := bytes.NewBuffer([]byte("test message"))
	url := fmt.Sprintf("http://%s/put?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, resp.StatusCode, 403)

	httpsAddr := nsqd.httpsListener.Addr().(*net.TCPAddr)
	cert, err := tls.LoadX509KeyPair("./test/certs/cert.pem", "./test/certs/key.pem")
	assert.Equal(t, err, nil)
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}

	buf = bytes.NewBuffer([]byte("test message"))
	url = fmt.Sprintf("https://%s/put?topic=%s", httpsAddr, topicName)
	resp, err = client.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(1))
}

func TestHTTPSRequireVerify(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	options := NewNSQDOptions()
	options.Verbose = true
	options.TLSCert = "./test/certs/server.pem"
	options.TLSKey = "./test/certs/server.key"
	options.TLSRootCAFile = "./test/certs/ca.pem"
	options.TLSClientAuthPolicy = "require-verify"
	_, httpAddr, nsqd := mustStartNSQD(options)
	httpsAddr := nsqd.httpsListener.Addr().(*net.TCPAddr)

	defer nsqd.Exit()

	topicName := "test_http_put_req_verf" + strconv.Itoa(int(time.Now().Unix()))
	topic := nsqd.GetTopic(topicName)

	// no cert
	buf := bytes.NewBuffer([]byte("test message"))
	url := fmt.Sprintf("http://%s/put?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/octet-stream", buf)
	assert.Equal(t, resp.StatusCode, 403)

	// unsigned cert
	cert, err := tls.LoadX509KeyPair("./test/certs/cert.pem", "./test/certs/key.pem")
	assert.Equal(t, err, nil)
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}

	buf = bytes.NewBuffer([]byte("test message"))
	url = fmt.Sprintf("https://%s/put?topic=%s", httpsAddr, topicName)
	resp, err = client.Post(url, "application/octet-stream", buf)
	assert.NotEqual(t, err, nil)

	// signed cert
	cert, err = tls.LoadX509KeyPair("./test/certs/client.pem", "./test/certs/client.key")
	assert.Equal(t, err, nil)
	tlsConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client = &http.Client{Transport: transport}

	buf = bytes.NewBuffer([]byte("test message"))
	url = fmt.Sprintf("https://%s/put?topic=%s", httpsAddr, topicName)
	resp, err = client.Post(url, "application/octet-stream", buf)
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "OK")

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, topic.Depth(), int64(1))
}

func TestHTTPDeprecatedTopicChannel(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_topic_channel" + strconv.Itoa(int(time.Now().Unix()))
	channelName := "ch"

	url := fmt.Sprintf("http://%s/create_topic?topic=%s", httpAddr, topicName)
	resp, err := http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	url = fmt.Sprintf("http://%s/create_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	topic, err := nsqd.GetExistingTopic(topicName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, topic, nil)

	channel, err := topic.GetExistingChannel(channelName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, channel, nil)

	url = fmt.Sprintf("http://%s/pause_topic?topic=%s", httpAddr, topicName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	assert.Equal(t, topic.IsPaused(), true)

	url = fmt.Sprintf("http://%s/unpause_topic?topic=%s", httpAddr, topicName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	assert.Equal(t, topic.IsPaused(), false)

	url = fmt.Sprintf("http://%s/pause_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	assert.Equal(t, channel.IsPaused(), true)

	url = fmt.Sprintf("http://%s/unpause_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	assert.Equal(t, channel.IsPaused(), false)

	url = fmt.Sprintf("http://%s/delete_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	_, err = topic.GetExistingChannel(channelName)
	assert.NotEqual(t, err, nil)

	url = fmt.Sprintf("http://%s/delete_topic?topic=%s", httpAddr, topicName)
	resp, err = http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), `{"status_code":200,"status_txt":"OK","data":null}`)

	_, err = nsqd.GetExistingTopic(topicName)
	assert.NotEqual(t, err, nil)
}

func TestHTTPTransitionTopicChannel(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	client := http.Client{}
	topicName := "test_http_topic_channel1" + strconv.Itoa(int(time.Now().Unix()))
	channelName := "ch1"

	url := fmt.Sprintf("http://%s/create_topic?topic=%s", httpAddr, topicName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err := client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	url = fmt.Sprintf("http://%s/create_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	topic, err := nsqd.GetExistingTopic(topicName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, topic, nil)

	channel, err := topic.GetExistingChannel(channelName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, channel, nil)

	url = fmt.Sprintf("http://%s/pause_topic?topic=%s", httpAddr, topicName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, topic.IsPaused(), true)

	url = fmt.Sprintf("http://%s/unpause_topic?topic=%s", httpAddr, topicName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, topic.IsPaused(), false)

	url = fmt.Sprintf("http://%s/pause_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, channel.IsPaused(), true)

	url = fmt.Sprintf("http://%s/unpause_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, channel.IsPaused(), false)

	url = fmt.Sprintf("http://%s/delete_channel?topic=%s&channel=%s", httpAddr, topicName, channelName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	_, err = topic.GetExistingChannel(channelName)
	assert.NotEqual(t, err, nil)

	url = fmt.Sprintf("http://%s/delete_topic?topic=%s", httpAddr, topicName)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err = client.Do(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	_, err = nsqd.GetExistingTopic(topicName)
	assert.NotEqual(t, err, nil)
}

func TestHTTPV1TopicChannel(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	_, httpAddr, nsqd := mustStartNSQD(NewNSQDOptions())
	defer nsqd.Exit()

	topicName := "test_http_topic_channel2" + strconv.Itoa(int(time.Now().Unix()))
	channelName := "ch2"

	url := fmt.Sprintf("http://%s/topic/create?topic=%s", httpAddr, topicName)
	resp, err := http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	url = fmt.Sprintf("http://%s/channel/create?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	topic, err := nsqd.GetExistingTopic(topicName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, topic, nil)

	channel, err := topic.GetExistingChannel(channelName)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, channel, nil)

	url = fmt.Sprintf("http://%s/topic/pause?topic=%s", httpAddr, topicName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, topic.IsPaused(), true)

	url = fmt.Sprintf("http://%s/topic/unpause?topic=%s", httpAddr, topicName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, topic.IsPaused(), false)

	url = fmt.Sprintf("http://%s/channel/pause?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, channel.IsPaused(), true)

	url = fmt.Sprintf("http://%s/channel/unpause?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	assert.Equal(t, channel.IsPaused(), false)

	url = fmt.Sprintf("http://%s/channel/delete?topic=%s&channel=%s", httpAddr, topicName, channelName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	_, err = topic.GetExistingChannel(channelName)
	assert.NotEqual(t, err, nil)

	url = fmt.Sprintf("http://%s/topic/delete?topic=%s", httpAddr, topicName)
	resp, err = http.Post(url, "application/json", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, 200)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, string(body), "")
	assert.Equal(t, resp.Header.Get("X-NSQ-Content-Type"), "nsq; version=1.0")

	_, err = nsqd.GetExistingTopic(topicName)
	assert.NotEqual(t, err, nil)
}

func BenchmarkHTTPput(b *testing.B) {
	var wg sync.WaitGroup
	b.StopTimer()
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	options := NewNSQDOptions()
	options.MemQueueSize = int64(b.N)
	_, httpAddr, nsqd := mustStartNSQD(options)
	msg := make([]byte, 256)
	topicName := "bench_http_put" + strconv.Itoa(int(time.Now().Unix()))
	url := fmt.Sprintf("http://%s/put?topic=%s", httpAddr, topicName)
	client := &http.Client{}
	b.SetBytes(int64(len(msg)))
	b.StartTimer()

	for j := 0; j < runtime.GOMAXPROCS(0); j++ {
		wg.Add(1)
		go func() {
			num := b.N / runtime.GOMAXPROCS(0)
			for i := 0; i < num; i += 1 {
				buf := bytes.NewBuffer(msg)
				req, _ := http.NewRequest("POST", url, buf)
				resp, err := client.Do(req)
				if err != nil {
					panic(err.Error())
				}
				body, _ := ioutil.ReadAll(resp.Body)
				if !bytes.Equal(body, []byte("OK")) {
					panic("bad response")
				}
				resp.Body.Close()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	b.StopTimer()
	nsqd.Exit()
}
