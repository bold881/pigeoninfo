package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WebservClient struct {
	pConn *websocket.Conn
}

type InConnClients struct {
	sync.RWMutex
	conns map[string]WebservClient
}

func (inConnClients *InConnClients) Init() {
	inConnClients.conns = make(map[string]WebservClient)
}

func (inConnClients *InConnClients) Check(ipport string) bool {
	inConnClients.Lock()
	exist := inConnClients.conns[ipport]
	inConnClients.Unlock()
	return (exist.pConn != nil)
}

func (inConnClients *InConnClients) Add(ipport string, conn *websocket.Conn) {
	if len(ipport) == 0 {
		return
	}
	inConnClients.Lock()
	if inConnClients.conns == nil {
		inConnClients.conns = make(map[string]WebservClient)
	}
	webServClient := WebservClient{conn}
	inConnClients.conns[ipport] = webServClient
	inConnClients.Unlock()
}

func (inConnClients *InConnClients) Del(ipport string) {
	inConnClients.Lock()
	if inConnClients.conns == nil {
		inConnClients.conns = make(map[string]WebservClient)
	}
	delete(inConnClients.conns, ipport)
	inConnClients.Unlock()
}

func (inConnClients *InConnClients) Clean() {
	for ipaddr, p := range inConnClients.conns {
		p.pConn.Close()
		delete(inConnClients.conns, ipaddr)
	}
}

type PageItemLite struct {
	Title string `json:"title"`
	Meta  string `json:"meta"`
}
