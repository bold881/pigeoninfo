package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var targetClass string = "news_main clearfix"

var crwedUrls CrawledURLs

type CommonInter interface {
	getSeedUrl() string
	isAllowed(szurl string) bool
	isDisallowed(szurl string) bool
}

type PageProcessor interface {
	pageProcess(szurl string, doc *goquery.Document, chPI chan PageItem)
	getCommonObj() *CommonObj
}

type CommonObj struct {
	url              string
	encode           string
	allowedDomain    []string
	disallowedDomain []string
}

type CQXinhuaObj struct {
	common CommonObj
}

type CQQQObj struct {
	common CommonObj
}

// check url suitable to crawl
func (p *CommonObj) checkUrl(szurl string) bool {
	if strings.Index(szurl, "http") == 0 {
		for _, vdis := range p.disallowedDomain {
			if strings.Contains(szurl, vdis) {
				return false
			}
		}

		for _, v := range p.allowedDomain {
			if strings.Contains(szurl, v) {
				return true
			}
		}
		return false
	}
	return false
}

// get item from
func (p *CQXinhuaObj) pageProcess(szurl string, doc *goquery.Document, chPI chan PageItem) {
	doc.Find(".news_main").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find(".article_title").Text()
		meta := s.Find(".time").Text()
		content := s.Find("p").Text()

		chPI <- PageItem{szurl, title, meta, content}
	})
}

func (p *CQQQObj) pageProcess(szurl string, doc *goquery.Document, chPI chan PageItem) {
	if doc == nil {
		return
	}
	doc.Find(".qq_article").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find(".hd").Find("h1").Text()
		if title == "" {
			log.Println("title empty" + szurl)
		}
		meta := s.Find(".a_Info").Text()
		if meta == "" {
			log.Println("meta empty" + szurl)
		}
		var content string
		s.Find(".Cnt-Main-Article-QQ").Find("p").Children().Remove().End().Each(func(j int, js *goquery.Selection) {
			if tmpTxt := js.Text(); tmpTxt != "" {
				content += js.Text() + "<br />"
			}
		})
		if content == "" {
			log.Println("content empty" + szurl)
		}

		chPI <- PageItem{szurl, title, meta, content}
	})
}

func (p *CQXinhuaObj) getCommonObj() *CommonObj {
	return &p.common
}

func (p *CQQQObj) getCommonObj() *CommonObj {
	return &p.common
}

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

// func CrawlGoQuery(szurl string, chPI chan PageItem, ch2Crawl chan string, isSeed bool) {

// 	if !isSeed {
// 		if crwedUrls.Check(szurl) {
// 			return
// 		}
// 	}

// 	crwedUrls.Add(szurl)

// 	doc, err := goquery.NewDocument(szurl)
// 	if err != nil {
// 		log.Println(err)
// 		crwedUrls.Del(szurl)
// 		return
// 	}

// 	// news article
// 	doc.Find(".news_main").Each(func(i int, s *goquery.Selection) {
// 		// For each item found, get the band and title
// 		title := s.Find(".article_title").Text()
// 		meta := s.Find(".time").Text()
// 		content := s.Find("p").Text()

// 		chPI <- PageItem{szurl, title, meta, content}
// 	})

// 	// all urls this page
// 	doc.Find("a").Each(func(i int, s *goquery.Selection) {
// 		href, _ := s.Attr("href")
// 		_, err := url.ParseRequestURI(href)
// 		if err == nil {
// 			if checkUrl(href) {
// 				if !crwedUrls.Check(href) {
// 					ch2Crawl <- href
// 				}
// 			}
// 		}
// 	})
// }

func GoScrapeRootOnly(pageProcItem PageProcessor, chPI chan PageItem, ch2Crawl chan string, isSeed bool) {

	if !isSeed {
		if crwedUrls.Check(pageProcItem.getCommonObj().url) {
			return
		}
	}

	crwedUrls.Add(pageProcItem.getCommonObj().url)

	var doc *goquery.Document
	var err error
	if pageProcItem.getCommonObj().encode == "utf-8" {
		doc, err = goquery.NewDocument(pageProcItem.getCommonObj().url)
	} else {
		res, err := http.Get(pageProcItem.getCommonObj().url)
		if err != nil {
			log.Println(err)
			crwedUrls.Del(pageProcItem.getCommonObj().url)
			return
		}
		defer res.Body.Close()

		utfBody, err := DecodeHTMLBody(res.Body, pageProcItem.getCommonObj().encode)
		if err != nil {
			log.Println(err)
			crwedUrls.Del(pageProcItem.getCommonObj().url)
			return
		}

		doc, err = goquery.NewDocumentFromReader(utfBody)
	}

	if err != nil {
		log.Println(err)
		crwedUrls.Del(pageProcItem.getCommonObj().url)
		return
	}

	// news article
	pageProcItem.pageProcess(pageProcItem.getCommonObj().url,
		doc, chPI)

	// all urls this page
	if isSeed {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			_, err := url.ParseRequestURI(href)
			if err == nil {
				if pageProcItem.getCommonObj().checkUrl(href) {
					if !crwedUrls.Check(href) {
						ch2Crawl <- href
					}
				}
			}
		})
	}
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// DecodeHTMLBody returns an decoding reader of the html Body for the specified `charset`
// If `charset` is empty, DecodeHTMLBody tries to guess the encoding from the content
func DecodeHTMLBody(body io.Reader, charset string) (io.Reader, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	return body, nil
}
