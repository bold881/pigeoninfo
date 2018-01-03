package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	//"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var inServiceClients InConnClients

func servEcho(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	inServiceClients.Add(r.RemoteAddr, c)

	// for {
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	for {
	// 		err = c.WriteMessage(mt, []byte("111"))
	// 		if err != nil {
	// 			log.Println("write:", err)
	// 			break
	// 		}
	// 		time.Sleep(time.Second * 5)
	// 	}
	// }
}

func broadcastLiteItem() {
	for {
		item := <-newsItemLiteChan
		for ipport, conn := range inServiceClients.conns {
			err := conn.pConn.WriteMessage(websocket.TextMessage, item)
			if err != nil {
				log.Println(err)
				inServiceClients.Del(ipport)
			}
		}
		//time.Sleep(time.Second * 5)
	}
}
