package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	session *mgo.Session
)

type News struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Title   string        `json:"title"`
	Meta    string        `json:"meta"`
	Content string        `json:"content"`
	Url     string        `json:"szurl"`
	Time    time.Time     `json:"sztime"`
}

func getTimeByStr(s string) (tearly time.Time, tlate time.Time, err error) {
	t1, err := time.Parse("2006-1-2", s)
	if err != nil {
		return
	} else {
		tearly = t1
	}

	s += " 23:59:59"
	t2, err := time.Parse("2006-1-2 15:04:05", s)
	if err != nil {
		return
	} else {
		tlate = t2
	}
	return
}

func GetNewsOfDay(t string) ([]News, error) {
	if session == nil {
		return nil, errors.New("database session not initialized")
	}
	cs := session.Copy()
	defer cs.Close()

	c := cs.DB("pigeoninfo").C("news")

	var results []News
	tearly, tlate, _ := getTimeByStr(t)
	m := bson.M{"time": bson.M{"$gte": tearly, "$lte": tlate}}
	err := c.Find(m).Sort("-time").All(&results)
	if err != nil {
		panic(err)
	}
	return results, nil
}

func GetNewsDetail(id string) (News, error) {
	var result News
	if session == nil {
		return result, errors.New("database session not initialized")
	}
	cs := session.Copy()
	defer cs.Close()

	c := cs.DB("pigeoninfo").C("news")

	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
	// m := bson.M{"_id": bson.ObjectIdHex{id}}
	// err := c.Find(m).One(&result)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func parseTime(s string) (t2 time.Time, err error) {
	t2, err = time.Parse("2006-01-02T15:04:05-07:00", s)
	return
}

func GetNewsOfLimit(ts string) ([]News, error) {
	if session == nil {
		return nil, errors.New("database session not initialized")
	}
	cs := session.Copy()
	defer cs.Close()

	c := cs.DB("pigeoninfo").C("news")

	var results []News
	//tearly, tlate, _ := getTimeByStr(t)

	var err error
	if ts == "" {
		iter := c.Find(nil).Sort("-time").Limit(25).Iter()
		err = iter.All(&results)
	} else {
		tlate, _ := parseTime(ts)
		m := bson.M{"time": bson.M{"$lte": tlate}}
		iter := c.Find(m).Sort("-time").Limit(25).Iter()
		err = iter.All(&results)
	}

	if err != nil {
		panic(err)
	}
	return results, nil
}
