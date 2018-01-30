package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"time"
)

var serverAddr = "0.0.0.0"
var reportPath = "http://0.0.0.0:4567/newsitem"

func save(ch2Save chan PageItem, pIs PageItems, s *mgo.Session) {
	for {
		pageItem := <-ch2Save
		if pageItem.content == "" {
			continue
		}
		if !MgoSave(s, pageItem) {
			fmt.Println(pageItem.title, pageItem.meta)
		}
		reportItem(&pageItem)
	}
}

// func add2crawl(ch2Crawl chan string, chPI chan PageItem) {
// 	for {
// 		url := <-ch2Crawl
// 		go CrawlGoQuery(url, chPI, ch2Crawl, false)
// 	}
// }

func stringInSlice(a string, list []PageProcessor) bool {
	for _, b := range list {
		if b.getCommonObj().url == a {
			return true
		}
	}
	return false
}

func initSeeds() []PageProcessor {
	var seedItems []PageProcessor
	var cqXinhuaObj CQXinhuaObj
	cqXinhuaObj.common.url = "http://www.cq.xinhuanet.com/"
	cqXinhuaObj.common.encode = "utf-8"
	cqXinhuaObj.common.allowedDomain = []string{
		"cq.xinhuanet.com",
	}
	cqXinhuaObj.common.disallowedDomain = []string{
		"big5.xinhuanet.com",
	}
	seedItems = append(seedItems, &cqXinhuaObj)

	var cqQQObj CQQQObj
	cqQQObj.common.url = "http://cq.qq.com"
	cqQQObj.common.encode = "gb2312"
	cqQQObj.common.allowedDomain = []string{
		"cq.qq.com",
	}
	seedItems = append(seedItems, &cqQQObj)
	return seedItems
}

func getProcItem(szurl string) (PageProcessor, bool) {
	if strings.Contains(szurl, "cq.qq.com") {
		var cqQQObj CQQQObj
		cqQQObj.common.url = szurl
		cqQQObj.common.encode = "gb2312"
		cqQQObj.common.allowedDomain = []string{
			"cq.qq.com",
		}
		return &cqQQObj, true
	}

	if strings.Contains(szurl, "cq.xinhuanet.com") {
		var cqXinhuaObj CQXinhuaObj
		cqXinhuaObj.common.url = szurl
		cqXinhuaObj.common.encode = "utf-8"
		cqXinhuaObj.common.allowedDomain = []string{
			"cq.xinhuanet.com",
		}
		cqXinhuaObj.common.disallowedDomain = []string{
			"big5.xinhuanet.com",
		}
		return &cqXinhuaObj, true
	}

	return nil, false
}

func main() {
	seedItems := initSeeds()

	var crawedItems PageItems
	crawedItems.Init()

	//var crawledUrls CrawledURLs
	crwedUrls.Init()

	session, err := mgo.Dial(serverAddr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	EnsureIndex(session)

	// Get all scraped URLs
	urls := GetCrawledUrls(session)

	for _, szurl := range urls {
		if !stringInSlice(szurl.Url, seedItems) {
			crwedUrls.Add(szurl.Url)
		}
	}

	ch2Crawl := make(chan string, 1000)
	chPageItem := make(chan PageItem, 1000)

	for _, pageProcItem := range seedItems {
		//go CrawlGoQuery(url, chPageItem, ch2Crawl, true)
		go GoScrapeRootOnly(pageProcItem, chPageItem, ch2Crawl, true)
	}

	go save(chPageItem, crawedItems, session)

	counter := 0
	for {
		if len(ch2Crawl) > 0 {
			url := <-ch2Crawl
			//go CrawlGoQuery(url, chPageItem, ch2Crawl, false)
			ppi, _ := getProcItem(url)
			if ppi != nil {
				go GoScrapeRootOnly(ppi, chPageItem, ch2Crawl, false)
			}

			time.Sleep(10 * time.Millisecond)
			continue
		} else if counter < 30 {
			counter++
			log.Println("sleep 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		if len(chPageItem) == 0 && len(ch2Crawl) == 0 {
			log.Println("rewind now...")
			//time.Sleep(10 * time.Minute)
			counter = 0
			for _, item := range seedItems {
				//go CrawlGoQuery(url, chPageItem, ch2Crawl, true)
				go GoScrapeRootOnly(item, chPageItem, ch2Crawl, true)
			}
		}
	}
}
