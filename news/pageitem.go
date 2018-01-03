package main

import (
	"sync"
)

type PageItem struct {
	url     string
	title   string
	meta    string
	content string
}

type PageItemLite struct {
	Title string `json:"title"`
	Meta  string `json:"meta"`
}

func (pi *PageItem) ToLite() PageItemLite {
	return PageItemLite{Title: pi.title, Meta: pi.meta}
}

// struct PageItems
// methods Find, Add
type PageItems struct {
	pageItems map[string]PageItem
}

func (pIs *PageItems) Find(url string) PageItem {
	if pIs.pageItems == nil {
		return PageItem{}
	}

	return pIs.pageItems[url]
}

func (pIs *PageItems) Add(url, title, meta, content string) {
	if pIs.pageItems == nil {
		pIs.pageItems = make(map[string]PageItem)
	}
	pIs.pageItems[url] = PageItem{url, title, meta, content}
}

func (pIs *PageItems) Init() {
	pIs.pageItems = make(map[string]PageItem)
}

// struct CrawledURLs,
// methods Check, Add
type CrawledURLs struct {
	sync.RWMutex
	crawled map[string]bool
}

func (cURLs *CrawledURLs) Init() {
	cURLs.crawled = make(map[string]bool)
}

func (cURLs *CrawledURLs) Check(url string) bool {
	cURLs.Lock()
	exist := cURLs.crawled[url]
	cURLs.Unlock()
	return exist
}

func (cURLs *CrawledURLs) Add(url string) {
	if len(url) == 0 {
		return
	}
	cURLs.Lock()
	if cURLs.crawled == nil {
		cURLs.crawled = make(map[string]bool)
	}
	cURLs.crawled[url] = true
	cURLs.Unlock()
}

func (cURLs *CrawledURLs) Del(url string) {
	cURLs.Lock()
	if cURLs.crawled == nil {
		cURLs.crawled = make(map[string]bool)
	}
	delete(cURLs.crawled, url)
	cURLs.Unlock()
}
