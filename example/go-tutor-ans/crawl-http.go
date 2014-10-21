//基于go totur的Web 爬虫， 实现HTTP爬虫
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (string, []string, error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, crawled map[string]struct{}, wait *sync.WaitGroup) {
	defer wait.Done()

	if depth <= 0 {
		return
	}
	// Don't fetch the same URL twice.
	if _, exists := crawled[url]; exists {
		fmt.Printf("already: %s\n", url)
		return
	}
	crawled[url] = struct{}{}
	//body, urls, err := fetcher.Fetch(url)
	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("found: %s %q\n", url, body)
	fmt.Printf("found: %s\n", url)

	for _, u := range urls {
		fmt.Println("crawling:", u)
		wait.Add(1)
		// Fetch URLs in parallel.
		go Crawl(u, depth-1, fetcher, crawled, wait)
	}
	//return
}

func main() {
	crawled := make(map[string]struct{})
	var wait sync.WaitGroup
	wait.Add(1)
	go Crawl("http://www.cyzly.sh.cn/amujap/", 2, fetcher, crawled, &wait)
	wait.Wait()
}

type HTTPFetcher struct {
	// TODO
}

var fetcher = &HTTPFetcher{}

func (f *HTTPFetcher) DoSomething(body string, urls []string) {
	// TODO
}

// 匹配 HTTP URL
var reg = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)

func (f *HTTPFetcher) Fetch(url string) (string, []string, error) {
	r, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer r.Body.Close()

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	body := string(content)
	urls := reg.FindAllString(body, -1)

	f.DoSomething(body, urls)

	return body, urls, nil
}
