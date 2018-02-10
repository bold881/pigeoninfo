package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"
)

const token string = "7173543247dff881"

func wxHandle(w http.ResponseWriter, req *http.Request) {

	parameters := req.URL.Query()
	if parameters == nil {
		return
	}

	timestamp := parameters.Get("timestamp")
	nonce := parameters.Get("nonce")
	echostr := parameters.Get("echostr")
	signature := parameters.Get("signature")

	var tmpArr []string
	tmpArr = append(tmpArr, token)
	tmpArr = append(tmpArr, timestamp)
	tmpArr = append(tmpArr, nonce)
	sort.Strings(tmpArr)
	var tmpStr string
	for _, ele := range tmpArr {
		tmpStr += ele
	}
	var retStr string
	tmpHashStr := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))
	if signature == tmpHashStr {
		retStr = echostr
	} else {
		retStr = "wx works!"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(retStr))
}
