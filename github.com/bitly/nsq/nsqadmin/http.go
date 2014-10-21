package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/nsq/nsqadmin/templates"
	"github.com/bitly/nsq/util"
	"github.com/bitly/nsq/util/lookupd"
	"github.com/bitly/nsq/util/semver"
)

var v1EndpointVersion *semver.Version

func init() {
	v1EndpointVersion, _ = semver.Parse("0.2.29-alpha")
}

// this is similar to httputil.NewSingleHostReverseProxy except it passes along basic auth
func NewSingleHostReverseProxy(target *url.URL, timeout time.Duration) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if target.User != nil {
			passwd, _ := target.User.Password()
			req.SetBasicAuth(target.User.Username(), passwd)
		}
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: util.NewDeadlineTransport(timeout),
	}
}

type httpServer struct {
	context  *Context
	counters map[string]map[string]int64
	proxy    *httputil.ReverseProxy
}

func NewHTTPServer(context *Context) *httpServer {
	var proxy *httputil.ReverseProxy

	templates.T.Funcs(template.FuncMap{
		"commafy":        util.Commafy,
		"nanotohuman":    util.NanoSecondToHuman,
		"floatToPercent": util.FloatToPercent,
		"percSuffix":     util.PercSuffix,
		"getNodeConsistencyClass": func(node *lookupd.Producer) string {
			if node.IsInconsistent(len(context.nsqadmin.options.NSQLookupdHTTPAddresses)) {
				return "btn-warning"
			}
			return ""
		},
	})

	templates.Parse()

	if context.nsqadmin.options.ProxyGraphite {
		url, err := url.Parse(context.nsqadmin.options.GraphiteURL)
		if err != nil {
			log.Fatalf("ERROR: failed to parse --graphite-url='%s' - %s",
				context.nsqadmin.options.GraphiteURL, err.Error())
		}
		proxy = NewSingleHostReverseProxy(url, 20*time.Second)
	}

	return &httpServer{
		context:  context,
		counters: make(map[string]map[string]int64),
		proxy:    proxy,
	}
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if strings.HasPrefix(req.URL.Path, "/node/") {
		s.nodeHandler(w, req)
		return
	} else if strings.HasPrefix(req.URL.Path, "/topic/") {
		s.topicHandler(w, req)
		return
	} else if strings.HasPrefix(req.URL.Path, "/static/") {
		if req.Method != "GET" {
			log.Printf("ERROR: invalid %s to GET only method", req.Method)
			http.Error(w, "INVALID_REQUEST", 500)
		} else {
			s.embeddedAssetHandler(w, req)
		}
		return
	}

	switch req.URL.Path {
	case "/":
		s.indexHandler(w, req)
	case "/ping":
		s.pingHandler(w, req)
	case "/nodes":
		s.nodesHandler(w, req)
	case "/tombstone_topic_producer":
		s.tombstoneTopicProducerHandler(w, req)
	case "/empty_topic":
		s.emptyTopicHandler(w, req)
	case "/delete_topic":
		s.deleteTopicHandler(w, req)
	case "/pause_topic":
		s.pauseTopicHandler(w, req)
	case "/unpause_topic":
		s.pauseTopicHandler(w, req)
	case "/delete_channel":
		s.deleteChannelHandler(w, req)
	case "/empty_channel":
		s.emptyChannelHandler(w, req)
	case "/pause_channel":
		s.pauseChannelHandler(w, req)
	case "/unpause_channel":
		s.pauseChannelHandler(w, req)
	case "/counter/data":
		s.counterDataHandler(w, req)
	case "/counter":
		s.counterHandler(w, req)
	case "/lookup":
		s.lookupHandler(w, req)
	case "/create_topic_channel":
		s.createTopicChannelHandler(w, req)
	case "/graphite_data":
		s.graphiteDataHandler(w, req)
	case "/render":
		if !s.context.nsqadmin.options.ProxyGraphite {
			http.NotFound(w, req)
			return
		}
		s.proxy.ServeHTTP(w, req)
	default:
		log.Printf("ERROR: 404 %s", req.URL.Path)
		http.NotFound(w, req)
	}
}

func (s *httpServer) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Length", "2")
	io.WriteString(w, "OK")
}

func (s *httpServer) embeddedAssetHandler(w http.ResponseWriter, req *http.Request) {
	var urlRegex = regexp.MustCompile(`^/static/(.+)$`)
	matches := urlRegex.FindStringSubmatch(req.URL.Path)
	if len(matches) == 0 {
		log.Printf("ERROR:  No embedded asset name for url - %s", req.URL.Path)
		http.NotFound(w, req)
		return
	}
	assetName := matches[1]
	log.Printf("INFO: Requesting embedded asset - %s", assetName)

	asset, error := templates.Asset(assetName)
	if error != nil {
		log.Printf("ERROR: embedded asset access - %s : %s", assetName, error)
		http.NotFound(w, req)
		return
	}
	assetLen := len(asset)

	if strings.HasSuffix(assetName, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(assetName, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", assetLen))
	w.Write(asset)
}

func (s *httpServer) indexHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	var topics []string
	if len(s.context.nsqadmin.options.NSQLookupdHTTPAddresses) != 0 {
		topics, _ = lookupd.GetLookupdTopics(s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
	} else {
		topics, _ = lookupd.GetNSQDTopics(s.context.nsqadmin.options.NSQDHTTPAddresses)
	}

	p := struct {
		Title        string
		GraphOptions *GraphOptions
		Topics       Topics
		Version      string
	}{
		Title:        "NSQ",
		GraphOptions: NewGraphOptions(w, req, reqParams, s.context),
		Topics:       TopicsFromStrings(topics),
		Version:      util.BINARY_VERSION,
	}
	err = templates.T.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

func (s *httpServer) topicHandler(w http.ResponseWriter, req *http.Request) {
	var urlRegex = regexp.MustCompile(`^/topic/(.*)$`)
	matches := urlRegex.FindStringSubmatch(req.URL.Path)
	if len(matches) == 0 {
		http.Error(w, "INVALID_TOPIC", 500)
		return
	}
	parts := strings.Split(matches[1], "/")
	topicName := parts[0]
	if !util.IsValidTopicName(topicName) {
		http.Error(w, "INVALID_TOPIC", 500)
		return
	}
	if len(parts) == 2 {
		channelName := parts[1]
		if !util.IsValidChannelName(channelName) {
			http.Error(w, "INVALID_CHANNEL", 500)
		} else {
			s.channelHandler(w, req, topicName, channelName)
		}
		return
	}

	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	producers := s.getProducers(topicName)
	topicStats, channelStats, _ := lookupd.GetNSQDStats(producers, topicName)

	globalTopicStats := &lookupd.TopicStats{HostAddress: "Total"}
	for _, t := range topicStats {
		globalTopicStats.Add(t)
	}

	hasE2eLatency := globalTopicStats.E2eProcessingLatency != nil &&
		len(globalTopicStats.E2eProcessingLatency.Percentiles) > 0

	var firstTopic *lookupd.TopicStats
	if len(topicStats) > 0 {
		firstTopic = topicStats[0]
	}

	p := struct {
		Title            string
		GraphOptions     *GraphOptions
		Version          string
		Topic            string
		TopicProducers   []string
		TopicStats       []*lookupd.TopicStats
		FirstTopic       *lookupd.TopicStats
		GlobalTopicStats *lookupd.TopicStats
		ChannelStats     map[string]*lookupd.ChannelStats
		HasE2eLatency    bool
	}{
		Title:            fmt.Sprintf("NSQ %s", topicName),
		GraphOptions:     NewGraphOptions(w, req, reqParams, s.context),
		Version:          util.BINARY_VERSION,
		Topic:            topicName,
		TopicProducers:   producers,
		TopicStats:       topicStats,
		FirstTopic:       firstTopic,
		GlobalTopicStats: globalTopicStats,
		ChannelStats:     channelStats,
		HasE2eLatency:    hasE2eLatency,
	}
	err = templates.T.ExecuteTemplate(w, "topic.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

func (s *httpServer) channelHandler(w http.ResponseWriter, req *http.Request, topicName string, channelName string) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	producers := s.getProducers(topicName)
	_, allChannelStats, _ := lookupd.GetNSQDStats(producers, topicName)
	channelStats := allChannelStats[channelName]

	hasE2eLatency := channelStats.E2eProcessingLatency != nil &&
		len(channelStats.E2eProcessingLatency.Percentiles) > 0

	var firstHost *lookupd.ChannelStats
	if len(channelStats.HostStats) > 0 {
		firstHost = channelStats.HostStats[0]
	}

	p := struct {
		Title          string
		GraphOptions   *GraphOptions
		Version        string
		Topic          string
		Channel        string
		TopicProducers []string
		ChannelStats   *lookupd.ChannelStats
		FirstHost      *lookupd.ChannelStats
		HasE2eLatency  bool
	}{
		Title:          fmt.Sprintf("NSQ %s / %s", topicName, channelName),
		GraphOptions:   NewGraphOptions(w, req, reqParams, s.context),
		Version:        util.BINARY_VERSION,
		Topic:          topicName,
		Channel:        channelName,
		TopicProducers: producers,
		ChannelStats:   channelStats,
		FirstHost:      firstHost,
		HasE2eLatency:  hasE2eLatency,
	}

	err = templates.T.ExecuteTemplate(w, "channel.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

func (s *httpServer) lookupHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	channels := make(map[string][]string)
	allTopics, _ := lookupd.GetLookupdTopics(s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
	for _, topicName := range allTopics {
		var producers []string
		producers, _ = lookupd.GetLookupdTopicProducers(topicName, s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
		if len(producers) == 0 {
			topicChannels, _ := lookupd.GetLookupdTopicChannels(topicName, s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
			channels[topicName] = topicChannels
		}
	}

	p := struct {
		Title        string
		GraphOptions *GraphOptions
		TopicMap     map[string][]string
		Lookupd      []string
		Version      string
	}{
		Title:        "NSQ Lookup",
		GraphOptions: NewGraphOptions(w, req, reqParams, s.context),
		TopicMap:     channels,
		Lookupd:      s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		Version:      util.BINARY_VERSION,
	}
	err = templates.T.ExecuteTemplate(w, "lookup.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

func (s *httpServer) createTopicChannelHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, err := reqParams.Get("topic")
	if err != nil || !util.IsValidTopicName(topicName) {
		http.Error(w, "INVALID_TOPIC", 500)
		return
	}

	channelName, err := reqParams.Get("channel")
	if err != nil || (len(channelName) > 0 && !util.IsValidChannelName(channelName)) {
		http.Error(w, "INVALID_CHANNEL", 500)
		return
	}

	for _, addr := range s.context.nsqadmin.options.NSQLookupdHTTPAddresses {
		nsqlookupdVersion, err := lookupd.GetVersion(addr)
		if err != nil {
			log.Printf("ERROR: failed to get nsqlookupd %s version - %s", addr, err)
		}

		uri := "create_topic"
		if !nsqlookupdVersion.Less(v1EndpointVersion) {
			uri = "topic/create"
		}

		endpoint := fmt.Sprintf("http://%s/%s?topic=%s", addr,
			uri, url.QueryEscape(topicName))
		log.Printf("LOOKUPD: querying %s", endpoint)
		_, err = util.APIRequestNegotiateV1("POST", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: lookupd %s - %s", endpoint, err)
			continue
		}

		if len(channelName) > 0 {
			uri := "create_channel"
			if !nsqlookupdVersion.Less(v1EndpointVersion) {
				uri = "channel/create"
			}
			endpoint := fmt.Sprintf("http://%s/%s?topic=%s&channel=%s",
				addr, uri,
				url.QueryEscape(topicName),
				url.QueryEscape(channelName))
			log.Printf("LOOKUPD: querying %s", endpoint)
			_, err := util.APIRequestNegotiateV1("POST", endpoint, nil)
			if err != nil {
				log.Printf("ERROR: lookupd %s - %s", endpoint, err.Error())
				continue
			}
		}
	}

	s.notifyAdminAction("create_topic", topicName, "", "", req)

	if len(channelName) > 0 {
		// TODO: we can remove this when we push new channel information from nsqlookupd -> nsqd
		producerAddrs, _ := lookupd.GetLookupdTopicProducers(topicName,
			s.context.nsqadmin.options.NSQLookupdHTTPAddresses)

		performVersionNegotiatedRequestsToNSQD(
			s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
			producerAddrs,
			"create_channel",
			"channel/create",
			fmt.Sprintf("topic=%s&channel=%s",
				url.QueryEscape(topicName), url.QueryEscape(channelName)))

		s.notifyAdminAction("create_channel", topicName, channelName, "", req)
	}

	http.Redirect(w, req, "/lookup", 302)
}

func (s *httpServer) tombstoneTopicProducerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, err := reqParams.Get("topic")
	if err != nil {
		http.Error(w, "MISSING_ARG_TOPIC", 500)
		return
	}

	node, err := reqParams.Get("node")
	if err != nil {
		http.Error(w, "MISSING_ARG_NODE", 500)
		return
	}

	rd, _ := reqParams.Get("rd")
	if !strings.HasPrefix(rd, "/") {
		rd = "/"
	}

	// tombstone the topic on all the lookupds
	for _, addr := range s.context.nsqadmin.options.NSQLookupdHTTPAddresses {
		nsqlookupdVersion, err := lookupd.GetVersion(addr)
		if err != nil {
			log.Printf("ERROR: failed to get nsqlookupd %s version - %s", addr, err)
		}

		uri := "tombstone_topic_producer"
		if !nsqlookupdVersion.Less(v1EndpointVersion) {
			uri = "topic/tombstone"
		}

		endpoint := fmt.Sprintf("http://%s/%s?topic=%s&node=%s",
			addr, uri,
			url.QueryEscape(topicName), url.QueryEscape(node))
		log.Printf("LOOKUPD: querying %s", endpoint)
		_, err = util.APIRequestNegotiateV1("POST", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: lookupd %s - %s", endpoint, err.Error())
		}
	}

	nsqdVersion, err := lookupd.GetVersion(node)
	if err != nil {
		log.Printf("ERROR: failed to get nsqd %s version - %s", node, err)
	}

	uri := "delete_topic"
	if !nsqdVersion.Less(v1EndpointVersion) {
		uri = "topic/delete"
	}

	// delete the topic on the producer
	endpoint := fmt.Sprintf("http://%s/%s?topic=%s", node,
		uri, url.QueryEscape(topicName))
	log.Printf("NSQD: querying %s", endpoint)
	_, err = util.APIRequestNegotiateV1("POST", endpoint, nil)
	if err != nil {
		log.Printf("ERROR: nsqd %s - %s", endpoint, err.Error())
	}

	s.notifyAdminAction("tombstone_topic_producer", topicName, "", node, req)

	http.Redirect(w, req, rd, 302)
}

func (s *httpServer) deleteTopicHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, err := reqParams.Get("topic")
	if err != nil {
		http.Error(w, "MISSING_ARG_TOPIC", 500)
		return
	}

	rd, _ := reqParams.Get("rd")
	if !strings.HasPrefix(rd, "/") {
		rd = "/"
	}

	// for topic removal, you need to get all the producers *first*
	producerAddrs := s.getProducers(topicName)

	// remove the topic from all the lookupds
	for _, addr := range s.context.nsqadmin.options.NSQLookupdHTTPAddresses {
		nsqlookupdVersion, err := lookupd.GetVersion(addr)
		if err != nil {
			log.Printf("ERROR: failed to get nsqlookupd %s version - %s", addr, err)
		}

		uri := "delete_topic"
		if !nsqlookupdVersion.Less(v1EndpointVersion) {
			uri = "topic/delete"
		}

		endpoint := fmt.Sprintf("http://%s/%s?topic=%s", addr, uri, topicName)
		log.Printf("LOOKUPD: querying %s", endpoint)
		_, err = util.APIRequestNegotiateV1("POST", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: lookupd %s - %s", endpoint, err.Error())
			continue
		}
	}

	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		"delete_topic",
		"topic/delete",
		fmt.Sprintf("topic=%s", url.QueryEscape(topicName)))

	s.notifyAdminAction("delete_topic", topicName, "", "", req)

	http.Redirect(w, req, rd, 302)
}

func (s *httpServer) deleteChannelHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, channelName, err := util.GetTopicChannelArgs(reqParams)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rd, _ := reqParams.Get("rd")
	if !strings.HasPrefix(rd, "/") {
		rd = fmt.Sprintf("/topic/%s", url.QueryEscape(topicName))
	}

	for _, addr := range s.context.nsqadmin.options.NSQLookupdHTTPAddresses {
		nsqlookupdVersion, err := lookupd.GetVersion(addr)
		if err != nil {
			log.Printf("ERROR: failed to get nsqlookupd %s version - %s", addr, err)
		}

		uri := "delete_channel"
		if !nsqlookupdVersion.Less(v1EndpointVersion) {
			uri = "channel/delete"
		}

		endpoint := fmt.Sprintf("http://%s/%s?topic=%s&channel=%s",
			addr, uri,
			url.QueryEscape(topicName),
			url.QueryEscape(channelName))
		log.Printf("LOOKUPD: querying %s", endpoint)
		_, err = util.APIRequestNegotiateV1("POST", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: lookupd %s - %s", endpoint, err.Error())
			continue
		}
	}

	producerAddrs := s.getProducers(topicName)
	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		"delete_channel",
		"channel/delete",
		fmt.Sprintf("topic=%s&channel=%s",
			url.QueryEscape(topicName), url.QueryEscape(channelName)))

	s.notifyAdminAction("delete_channel", topicName, channelName, "", req)

	http.Redirect(w, req, rd, 302)
}

func (s *httpServer) emptyTopicHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, err := reqParams.Get("topic")
	if err != nil {
		http.Error(w, "MISSING_ARG_TOPIC", 500)
		return
	}

	producerAddrs := s.getProducers(topicName)
	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		"empty_topic",
		"topic/empty",
		fmt.Sprintf("topic=%s",
			url.QueryEscape(topicName)))

	s.notifyAdminAction("empty_topic", topicName, "", "", req)

	http.Redirect(w, req, fmt.Sprintf("/topic/%s", url.QueryEscape(topicName)), 302)
}

func (s *httpServer) pauseTopicHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, err := reqParams.Get("topic")
	if err != nil {
		http.Error(w, "MISSING_ARG_TOPIC", 500)
		return
	}

	verb := "pause"
	if strings.Contains(req.URL.Path, "unpause") {
		verb = "unpause"
	}

	producerAddrs := s.getProducers(topicName)
	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		verb+"_topic",
		"topic/"+verb,
		fmt.Sprintf("topic=%s",
			url.QueryEscape(topicName)))

	s.notifyAdminAction(verb+"_topic", topicName, "", "", req)

	http.Redirect(w, req, fmt.Sprintf("/topic/%s", url.QueryEscape(topicName)), 302)
}

func (s *httpServer) emptyChannelHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, channelName, err := util.GetTopicChannelArgs(reqParams)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	producerAddrs := s.getProducers(topicName)
	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		"empty_channel",
		"channel/empty",
		fmt.Sprintf("topic=%s&channel=%s",
			url.QueryEscape(topicName), url.QueryEscape(channelName)))

	s.notifyAdminAction("empty_channel", topicName, channelName, "", req)

	http.Redirect(w, req, fmt.Sprintf("/topic/%s/%s",
		url.QueryEscape(topicName), url.QueryEscape(channelName)), 302)
}

func (s *httpServer) pauseChannelHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("ERROR: invalid %s to POST only method", req.Method)
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	reqParams := &util.PostParams{req}

	topicName, channelName, err := util.GetTopicChannelArgs(reqParams)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	verb := "pause"
	if strings.Contains(req.URL.Path, "unpause") {
		verb = "unpause"
	}

	producerAddrs := s.getProducers(topicName)
	performVersionNegotiatedRequestsToNSQD(
		s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
		producerAddrs,
		verb+"_channel",
		"channel/"+verb,
		fmt.Sprintf("topic=%s&channel=%s",
			url.QueryEscape(topicName), url.QueryEscape(channelName)))

	s.notifyAdminAction(verb+"_channel", topicName, channelName, "", req)

	http.Redirect(w, req, fmt.Sprintf("/topic/%s/%s", url.QueryEscape(topicName), url.QueryEscape(channelName)), 302)
}

func (s *httpServer) nodeHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	var urlRegex = regexp.MustCompile(`^/node/(.*)$`)
	matches := urlRegex.FindStringSubmatch(req.URL.Path)
	if len(matches) == 0 {
		http.Error(w, "INVALID_NODE", 500)
		return
	}
	parts := strings.Split(matches[1], "/")
	node := parts[0]

	found := false
	for _, n := range s.context.nsqadmin.options.NSQDHTTPAddresses {
		if node == n {
			found = true
			break
		}
	}
	producers, _ := lookupd.GetLookupdProducers(s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
	for _, p := range producers {
		if node == fmt.Sprintf("%s:%d", p.BroadcastAddress, p.HttpPort) {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "INVALID_NODE", 500)
		return
	}

	topicStats, channelStats, _ := lookupd.GetNSQDStats([]string{node}, "")

	numClients := int64(0)
	numMessages := int64(0)
	for _, ts := range topicStats {
		for _, cs := range ts.Channels {
			numClients += int64(len(cs.Clients))
		}
		numMessages += ts.MessageCount
	}

	p := struct {
		Title        string
		Version      string
		GraphOptions *GraphOptions
		Node         Node
		TopicStats   []*lookupd.TopicStats
		ChannelStats map[string]*lookupd.ChannelStats
		NumMessages  int64
		NumClients   int64
	}{
		Title:        "NSQ Node - " + node,
		Version:      util.BINARY_VERSION,
		GraphOptions: NewGraphOptions(w, req, reqParams, s.context),
		Node:         Node(node),
		TopicStats:   topicStats,
		ChannelStats: channelStats,
		NumMessages:  numMessages,
		NumClients:   numClients,
	}
	err = templates.T.ExecuteTemplate(w, "node.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

func (s *httpServer) nodesHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	producers, _ := lookupd.GetLookupdProducers(s.context.nsqadmin.options.NSQLookupdHTTPAddresses)

	p := struct {
		Title        string
		Version      string
		GraphOptions *GraphOptions
		Producers    []*lookupd.Producer
		Lookupd      []string
	}{
		Title:        "NSQ Nodes",
		Version:      util.BINARY_VERSION,
		GraphOptions: NewGraphOptions(w, req, reqParams, s.context),
		Producers:    producers,
		Lookupd:      s.context.nsqadmin.options.NSQLookupdHTTPAddresses,
	}
	err = templates.T.ExecuteTemplate(w, "nodes.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

type counterTarget struct{}

func (c counterTarget) Target(key string) ([]string, string) {
	return []string{fmt.Sprintf("sumSeries(%%stopic.*.channel.*.%s)", key)}, "green"
}

func (c counterTarget) Host() string {
	return "*"
}

func (s *httpServer) counterHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}
	p := struct {
		Title        string
		Version      string
		GraphOptions *GraphOptions
		Target       counterTarget
	}{
		Title:        "NSQ Message Counts",
		Version:      util.BINARY_VERSION,
		GraphOptions: NewGraphOptions(w, req, reqParams, s.context),
		Target:       counterTarget{},
	}
	err = templates.T.ExecuteTemplate(w, "counter.html", p)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(w, "Template Error", 500)
	}
}

// this endpoint works by giving out an ID that maps to a stats dictionary
// The initial request is the number of messages processed since each nsqd started up.
// Subsequent requsts pass that ID and get an updated delta based on each individual channel/nsqd message count
// That ID must be re-requested or it will be expired.
func (s *httpServer) counterDataHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		util.ApiResponse(w, 500, "INVALID_REQUEST", nil)
		return
	}

	statsID, _ := reqParams.Get("id")
	now := time.Now()
	if statsID == "" {
		// make a new one
		statsID = fmt.Sprintf("%d.%d", now.Unix(), now.UnixNano())
	}

	stats, ok := s.counters[statsID]
	if !ok {
		stats = make(map[string]int64)
	}
	newStats := make(map[string]int64)
	newStats["time"] = now.Unix()

	producers, _ := lookupd.GetLookupdProducers(s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
	addresses := make([]string, len(producers))
	for i, p := range producers {
		addresses[i] = p.HTTPAddress()
	}
	_, channelStats, _ := lookupd.GetNSQDStats(addresses, "")

	var newMessages int64
	var totalMessages int64
	for _, channelStats := range channelStats {
		for _, hostChannelStats := range channelStats.HostStats {
			key := fmt.Sprintf("%s:%s:%s", channelStats.TopicName, channelStats.ChannelName, hostChannelStats.HostAddress)
			d, ok := stats[key]
			if ok && d <= hostChannelStats.MessageCount {
				newMessages += (hostChannelStats.MessageCount - d)
			}
			totalMessages += hostChannelStats.MessageCount
			newStats[key] = hostChannelStats.MessageCount
		}
	}
	s.counters[statsID] = newStats

	data := make(map[string]interface{})
	data["new_messages"] = newMessages
	data["total_messages"] = totalMessages
	data["id"] = statsID
	util.ApiResponse(w, 200, "OK", data)
}

func (s *httpServer) graphiteDataHandler(w http.ResponseWriter, req *http.Request) {
	reqParams, err := util.NewReqParams(req)
	if err != nil {
		log.Printf("ERROR: failed to parse request params - %s", err.Error())
		http.Error(w, "INVALID_REQUEST", 500)
		return
	}

	metric, err := reqParams.Get("metric")
	if err != nil {
		log.Printf("ERROR: missing metric param - %s", err.Error())
		http.Error(w, "MISSING_METRIC_PARAM", 500)
		return
	}

	target, err := reqParams.Get("target")
	if err != nil {
		log.Printf("ERROR: missing target param - %s", err.Error())
		http.Error(w, "MISSING_TARGET_PARAM", 500)
		return
	}

	var queryFunc func(string) string
	var formatJsonResponseFunc func([]byte) ([]byte, error)

	switch metric {
	case "rate":
		queryFunc = rateQuery
		formatJsonResponseFunc = parseRateResponse
	default:
		log.Printf("ERROR: unknown metric value %s", metric)
		http.Error(w, "INVALID_METRIC_PARAM", 500)
		return
	}

	query := queryFunc(target)
	url := s.context.nsqadmin.options.GraphiteURL + query
	log.Printf("GRAPHITE: %s", url)
	response, err := graphiteGet(url)
	if err != nil {
		log.Printf("ERROR: graphite request failed %s", err.Error())
		http.Error(w, "GRAPHITE_FAILED", 500)
		return
	}

	resp, err := formatJsonResponseFunc(response)
	if err != nil {
		log.Printf("ERROR: response formating failed - %s", err.Error())
		http.Error(w, "INVALID_GRAPHITE_RESPONSE", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return
}

func graphiteGet(request_url string) ([]byte, error) {
	response, err := http.Get(request_url)

	var contents []byte

	if err != nil {
		log.Printf("ERROR: GET request to graphite failed %s", err)
		return nil, err
	}

	defer response.Body.Close()
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: reading GET body failed %s", err)
		return nil, err
	}
	return contents, nil
}

func (s *httpServer) getProducers(topicName string) []string {
	var producers []string
	if len(s.context.nsqadmin.options.NSQLookupdHTTPAddresses) != 0 {
		producers, _ = lookupd.GetLookupdTopicProducers(topicName, s.context.nsqadmin.options.NSQLookupdHTTPAddresses)
	} else {
		producers, _ = lookupd.GetNSQDTopicProducers(topicName, s.context.nsqadmin.options.NSQDHTTPAddresses)
	}
	return producers
}

func producerSearch(producers []*lookupd.Producer, needle string) *lookupd.Producer {
	for _, producer := range producers {
		addr := net.JoinHostPort(producer.BroadcastAddress, strconv.Itoa(producer.HttpPort))
		if needle == addr {
			return producer
		}
	}
	return nil
}

func performVersionNegotiatedRequestsToNSQD(nsqlookupdAddrs []string, nsqdAddrs []string,
	deprecatedURI string, v1URI string, queryString string) {
	var err error
	// get producer structs in one set of up-front requests
	// so we can negotiate versions
	//
	// (this returns an empty list if there are no nsqlookupd configured)
	producers, _ := lookupd.GetLookupdProducers(nsqlookupdAddrs)

	for _, addr := range nsqdAddrs {
		var nodeVer *semver.Version

		uri := deprecatedURI
		producer := producerSearch(producers, addr)
		if producer != nil {
			nodeVer = producer.VersionObj
		} else {
			// we couldn't find the node in our list
			// so ask it for a version directly
			nodeVer, err = lookupd.GetVersion(addr)
			if err != nil {
				log.Printf("ERROR: failed to get nsqd %s version - %s", addr, err)
			}
		}

		if nodeVer != nil && !nodeVer.Less(v1EndpointVersion) {
			uri = v1URI
		}

		endpoint := fmt.Sprintf("http://%s/%s?%s", addr, uri, queryString)
		log.Printf("NSQD: querying %s", endpoint)
		_, err := util.APIRequestNegotiateV1("POST", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: nsqd %s - %s", endpoint, err.Error())
			continue
		}
	}
}
