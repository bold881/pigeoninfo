package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

type LogSet struct {
	sync.RWMutex
	logs []string
}

var (
	runninglog LogSet
)

// send logs by eamil
func logProcess() {
	for {
		tnow := time.Now()
		a := fmt.Sprintf("%d-%d-%d 23:30:00", tnow.Year(), tnow.Month(), tnow.Day())
		t, err := time.ParseInLocation("2006-1-2 15:04:05", a, tnow.Location())
		if err != nil {
			log.Println(err)
		}
		if !(tnow.Hour() > 23 && tnow.Minute() > 30) {

			log.Println(t.Sub(tnow))
			time.Sleep(t.Sub(tnow))

			msg := getEmailMsg(runninglog.logs)
			if msg != nil {
				mailsend(msg)
			}
		}
		tthen := time.Now()
		if tthen.Day() == tnow.Day() && tthen.Hour() == tnow.Hour() {
			tnext := t.AddDate(0, 0, 1)
			time.Sleep(tnext.Sub(tthen))
		}
	}
}

// send eamil
func mailsend(msg []byte) bool {
	c, err := smtp.Dial("mail.yonyou.com:25")
	if err != nil {
		log.Println(err)
		return false
	}

	cfg := &tls.Config{ServerName: "", InsecureSkipVerify: true}
	c.StartTLS(cfg)

	auth := smtp.PlainAuth("", "", "", "")
	c.Auth(auth)

	if err := c.Mail(""); err != nil {
		log.Println(err)
		return false
	}
	if err := c.Rcpt(""); err != nil {
		log.Println(err)
		return false
	}

	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = wc.Write(msg)
	if err != nil {
		log.Println(err)
		wc.Close()
		return false
	}
	err = wc.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	err = c.Quit()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func getEmailMsg(logs []string) []byte {
	if len(logs) < 1 {
		return nil
	}
	tnow := time.Now()
	title := fmt.Sprintf("%d-%d-%d RUNNING REPORT", tnow.Year(), tnow.Month(), tnow.Day())
	msg := ("To: 970778418@qq.com\r\n" +
		"Subject: " + title + "\r\n" +
		"\r\n")

	for _, s := range logs {
		msg += s + "\r\n"
	}

	return []byte(msg)
}
