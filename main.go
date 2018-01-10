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

// Config 타입은 config.json의 값들을 받기위한 구조체
type Config struct {
	PortalCookie  string `json:"portal_cookie"`
	FacebookToken string `json:"facebook_token"`
}

// config 변수는 전역변수로 사용
var config Config

// 초기화 - config 파일을 읽어서 전역변수에 저장
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
}

// 메인함수
func main() {
	// ctrl + c 눌렀을 때, ajou-noticer off를 알리는 부분
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("[ajou-noticer] ajou-noticer off")
		os.Exit(1)
	}()

	// checker를 선언하고, 5분마다 확인하도록 하는 부분
	checker := Checker{}
	fmt.Println("[ajou-noticer] ajou-noticer on")
	for {
		checker.check()
		time.Sleep(5 * time.Minute)
	}
}
