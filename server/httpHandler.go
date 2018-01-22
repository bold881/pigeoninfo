package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// get news of specific day
func newsofDay(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	bodyStr := string(bodyBytes)
	log.Println(bodyStr)
	runninglog.logs = append(runninglog.logs, bodyStr)
	news, err := GetNewsOfDay(bodyStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	retb, _ := json.Marshal(news)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-*, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Request-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Write(retb)
}

func newsDetail(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	bodyStr := string(bodyBytes)
	log.Println(bodyStr)
	runninglog.logs = append(runninglog.logs, bodyStr)
	news, err := GetNewsDetail(bodyStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	retb, _ := json.Marshal(news)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-*, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Request-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Write(retb)
}

// notify incoming news item
func newsItem(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err == nil {
		//bodyStr := string(bodyBytes)
		//log.Println(bodyStr)
		//item := PageItemLite{}
		//json.Unmarshal(bodyBytes, &item)
		//log.Println(item)
		newsItemLiteChan <- bodyBytes
	} else {
		log.Println(err)
		runninglog.logs = append(runninglog.logs, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
