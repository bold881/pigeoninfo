package main

import (
	"fmt"
	"time"
)

func pinger(c chan string) {
	for {
		c <- "ping"
	}
}

func pong(c chan string) {
	for {
		c <- "pong"
	}
}

func printer(c chan string) {
	for {
		fmt.Println(<-c)
		time.Sleep(time.Second)
	}
}

func Test() {
	c := make(chan string)
	go pinger(c)
	go printer(c)
	go pong(c)

	var input string
	fmt.Scanln(&input)
}
