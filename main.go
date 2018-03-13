package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
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
		config.FacebookToken = ""
		config.PortalCookie = ""
		return
	}
	err := json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
}

// 메인함수
func main() {
	fmt.Println("[ajou-noticer] ajou-noticer on")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// ctrl + c 눌렀을 때, ajou-noticer off를 알리는 부분
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("[ajou-noticer] ajou-noticer off")
		os.Exit(1)
	}()

	wg.Add(1)
	go func() {
		// checker를 선언하고, 5분마다 확인하도록 하는 부분
		checker := Checker{}
		fmt.Println("[ajou-noticer] checker on")
		for {
			checker.check()
			time.Sleep(5 * time.Minute)
		}
	}()

	wg.Add(1)
	go func() {
		fmt.Println("[ajou-noticer] server on")
		StartServer("51234")
	}()
	wg.Wait()
}
