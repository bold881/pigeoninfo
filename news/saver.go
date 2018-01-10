package main

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"hash/fnv"
	"log"
	"regexp"
	"strings"
	"time"
)

type News struct {
	//ID      bson.ObjectId `bson:"_id,omitempty"`
	Title   string
	Meta    string
	Content string
	Url     string
	Time    time.Time
}

var (
	session    mgo.Session
	timeregexp *regexp.Regexp
)

func init() {
	timeregexp = regexp.MustCompile(`\d{4}-(\d{2}|\d)-(\d{2}|\d) ([01]?[0-9]|2[0-3]):[0-5][0-9]`)
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func initSession(session *mgo.Session) bool {
	if session == nil {
		session, err := mgo.Dial("101.200.47.113")
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		return true
	}
	return true
}

func Clean(session *mgo.Session) {
	if session != nil {
		session.Close()
	}
}

func string2Time(rawstr string) time.Time {
	if tmpstr := timeregexp.FindString(rawstr); tmpstr != "" {
		t, err := time.Parse("2006-01-02 15:04", tmpstr)
		if err == nil {
			return t
		} else {
			return time.Now()
		}
	}

	var strarr []string
	if strings.Contains(rawstr, "来源") {
		strarr = strings.Split(rawstr, "来源")
	} else if strings.Contains(rawstr, "來源") {
		strarr = strings.Split(rawstr, "來源")
	} else {
		return time.Now()
	}
	rawstr = strarr[0]
	rawstr = strings.TrimSpace(rawstr)
	rawstr = strings.Replace(rawstr, "年", "-", -1)
	rawstr = strings.Replace(rawstr, "月", "-", -1)
	rawstr = strings.Replace(rawstr, "日", "", -1)
	t, err := time.Parse("2006-01-02 15:04", rawstr)
	if err == nil {
		return t
	} else {
		return time.Now()
	}
}

func transItem(pItem PageItem) News {
	t := string2Time(pItem.meta)
	//b := hash(pItem.url)
	//var str string = string(b) + string(b) + string(b) + string(b)

	return News{
		//bson.ObjectId(str),
		pItem.title,
		pItem.meta,
		pItem.content,
		pItem.url,
		t}
}

func MgoSave(session *mgo.Session, pItem PageItem) bool {

	s := session.Copy()
	if s != nil {
		defer s.Close()
	}

	c := s.DB("pigeoninfo").C("news")

	news := transItem(pItem)

	err := c.Insert(&news)
	if err != nil {
		log.Print(err)
		return false
	}

	return true
	// result := Person{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Phone:", result.Phone)
}

func EnsureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("pigeoninfo").C("news")

	index := mgo.Index{
		Key:        []string{"url"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
