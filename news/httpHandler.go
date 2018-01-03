package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func reportItem(pi *PageItem) {
	piLite := pi.ToLite()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(piLite)
	log.Println(piLite)
	log.Println(b)

	contentType := "application/json; charset=UTF-8"
	resp, err := http.Post(reportPath, contentType, b)
	if err == nil {
		resp.Body.Close()
	} else {
		log.Println(err)
	}
}
