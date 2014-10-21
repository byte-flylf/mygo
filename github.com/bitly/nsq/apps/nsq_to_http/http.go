package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/bitly/nsq/util"
)

var httpclient *http.Client
var userAgent string

func init() {
	httpclient = &http.Client{Transport: util.NewDeadlineTransport(*httpTimeout)}
	userAgent = fmt.Sprintf("nsq_to_http v%s", util.BINARY_VERSION)
}

func HttpGet(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return httpclient.Do(req)
}

func HttpPost(endpoint string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", *contentType)
	return httpclient.Do(req)
}
