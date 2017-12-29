package main

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetCrawledUrls(s *mgo.Session) []News {
	session := s.Copy()
	defer session.Close()

	c := session.DB("pigeoninfo").C("news")

	var results []News
	m := []bson.M{
		{"$project": bson.M{"url": 1, "_id": 0}},
	}
	err := c.Pipe(m).All(&results)
	if err != nil {
		panic(err)
	}
	//fmt.Println(results[0])
	return results
}
