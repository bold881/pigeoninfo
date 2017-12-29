package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var MGOADDR = "101.200.47.113"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MyHandler struct {
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr + " " + req.URL.Path)
	if req.URL.Path == "/" {
		io.WriteString(w, "hello world!\n")
	} else if req.URL.Path == "/newsofday" {
		if req.Method == "POST" {
			newsofDay(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	} else if req.URL.Path == "/echo" {
		echo(w, req)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func newsofDay(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	bodyStr := string(bodyBytes)
	log.Println(bodyStr)
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

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for {
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
			time.Sleep(time.Second * 5)
		}

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
