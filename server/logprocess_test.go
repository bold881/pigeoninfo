package main

import (
	"testing"
)

func TestMailsend(t *testing.T) {
	m := []string{
		"hello world",
		"this from daifenga",
	}

	msg := getEmailMsg(m)
	if msg != nil {
		if !mailsend(msg) {
			t.Error("send mail fail")
		}
	} else {
		t.Error("get email msg fail")
	}
}
