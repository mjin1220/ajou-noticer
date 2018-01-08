package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	PortalCookie  string `json:"portal_cookie"`
	FacebookToken string `json:"facebook_token"`
}

var config Config

func init() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	err := json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	// 초기화 함수
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("[ajou-noticer] ajou-noticer off")
		os.Exit(1)
	}()

	checker := Checker{}
	fmt.Println("[ajou-noticer] ajou-noticer on")
	for {
		checker.check()
		time.Sleep(5 * time.Minute)
	}
}
