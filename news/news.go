package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

var serverAddr = "10.115.0.29"
var reportPath = "http://10.115.0.134:4567/newsitem"

func save(ch2Save chan PageItem, pIs PageItems, s *mgo.Session) {
	for {
		pageItem := <-ch2Save
		if !MgoSave(s, pageItem) {
			fmt.Println(pageItem.title, pageItem.meta)
		}
		reportItem(&pageItem)
	}
}

func add2crawl(ch2Crawl chan string, chPI chan PageItem) {
	for {
		url := <-ch2Crawl
		go CrawlGoQuery(url, chPI, ch2Crawl, false)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	seedUrls := []string{
		"http://www.cq.xinhuanet.com/",
	}

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
		if !stringInSlice(szurl.Url, seedUrls) {
			crwedUrls.Add(szurl.Url)
		}
	}

	ch2Crawl := make(chan string, 1000)
	chPageItem := make(chan PageItem, 1000)

	for _, url := range seedUrls {
		//go CrawlGoQuery(url, chPageItem, ch2Crawl, true)
		go GoScrapeRootOnly(url, chPageItem, ch2Crawl, true)
	}

	go save(chPageItem, crawedItems, session)

	counter := 0
	for {
		if len(ch2Crawl) > 0 {
			url := <-ch2Crawl
			//go CrawlGoQuery(url, chPageItem, ch2Crawl, false)
			go GoScrapeRootOnly(url, chPageItem, ch2Crawl, false)
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
			for _, url := range seedUrls {
				//go CrawlGoQuery(url, chPageItem, ch2Crawl, true)
				go GoScrapeRootOnly(url, chPageItem, ch2Crawl, true)
			}
		}
	}
}
