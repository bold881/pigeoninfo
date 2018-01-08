package main

import (
	"gopkg.in/mgo.v2"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	MGOADDR = "101.200.47.113"
	//MGOADDR          = "10.115.0.29"
	newsItemLiteChan = make(chan []byte, 1000)
)

type MyHandler struct {
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr + " " + req.URL.Path)
	runninglog.logs = append(runninglog.logs, req.RemoteAddr+" "+req.URL.Path)
	if req.URL.Path == "/" {
		io.WriteString(w, "hello world!\n")
	} else if req.URL.Path == "/newsofday" {
		if req.Method == "POST" {
			newsofDay(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	} else if req.URL.Path == "/newsitem" {
		if req.Method == "POST" {
			newsItem(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	} else if req.URL.Path == "/echo" {
		servEcho(w, req)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	var err error
	session, err = mgo.Dial(MGOADDR)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	defer inServiceClients.Clean()

	go broadcastLiteItem()

	go logProcess()

	var handler MyHandler
	s := &http.Server{
		Addr:           ":4567",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
