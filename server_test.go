package main

import (
	"net/http"
	"testing"
)

func TestStartServer(t *testing.T) {
	http.Get("http://127.0.0.1:51234/release/AjouNoticer.crt")
	http.Get("http://127.0.0.1:51234/release/AjouNoticer.crt")
}
