package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	checker := Checker{}

	// checker.check()
	checker.sendMessage(checker.makeMessage())
}
