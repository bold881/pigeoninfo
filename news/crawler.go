package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	//"io/ioutil"
	"log"
	//"net/http"
	"net/url"
	"strings"
)

var allowedDomain string = "cq.xinhuanet.com"
var disallowedDomain string = "big5.xinhuanet.com"
var targetClass string = "news_main clearfix"

var crwedUrls CrawledURLs

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func getDivText(t html.Token, targetClass string) bool {
	for _, a := range t.Attr {
		if a.Key == "class" && a.Val == targetClass {
			return true
		}
	}
	return false
}

func checkUrl(szurl string) bool {
	if strings.Index(szurl, "http") == 0 {
		if strings.Contains(szurl, disallowedDomain) {
			return false
		}
		return strings.Contains(szurl, allowedDomain)
	}
	return false
}

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://www.cq.xinhuanet.com/2017-12/10/c_1122086150.htm")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".news_main").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find(".article_title").Text()
		info := s.Find(".time").Text()
		content := s.Find("p").Text()

		fmt.Println(title, info, content)
	})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		fmt.Printf("%d, %s\n", i, href)
	})
}

func CrawlGoQuery(szurl string, chPI chan PageItem, ch2Crawl chan string, isSeed bool) {

	if !isSeed {
		if crwedUrls.Check(szurl) {
			return
		}
	}

	crwedUrls.Add(szurl)

	doc, err := goquery.NewDocument(szurl)
	if err != nil {
		log.Println(err)
		crwedUrls.Del(szurl)
		return
	}

	// news article
	doc.Find(".news_main").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find(".article_title").Text()
		meta := s.Find(".time").Text()
		content := s.Find("p").Text()

		chPI <- PageItem{szurl, title, meta, content}
	})

	// all urls this page
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		_, err := url.ParseRequestURI(href)
		if err == nil {
			if checkUrl(href) {
				if !crwedUrls.Check(href) {
					ch2Crawl <- href
				}
			}
		}
	})
}

func GoScrapeRootOnly(szurl string, chPI chan PageItem, ch2Crawl chan string, isSeed bool) {

	if !isSeed {
		if crwedUrls.Check(szurl) {
			return
		}
	}

	crwedUrls.Add(szurl)

	doc, err := goquery.NewDocument(szurl)
	if err != nil {
		log.Println(err)
		crwedUrls.Del(szurl)
		return
	}

	// news article
	doc.Find(".news_main").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find(".article_title").Text()
		meta := s.Find(".time").Text()
		content := s.Find("p").Text()

		chPI <- PageItem{szurl, title, meta, content}
	})

	// all urls this page
	if isSeed {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			_, err := url.ParseRequestURI(href)
			if err == nil {
				if checkUrl(href) {
					if !crwedUrls.Check(href) {
						ch2Crawl <- href
					}
				}
			}
		})
	}
}
