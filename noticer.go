package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Noticer struct {
}

func (noticer Noticer) Notify(notices Notices) {
	for i := len(notices) - 1; i >= 0; i-- {
		noticer.sendMessage(noticer.makeMessage(notices[i]))
		fmt.Printf("[%s] %s\n", "New", notices[i].Title)
	}
}

func (noticer Noticer) makeMessage(notice Notice) (messageCreativeID string) {
	buttons := []Button{{"web_url", notice.URL, "자세히 보기"}}
	elements := []Element{{notice.Title, notice.Department, buttons}}
	payload := Payload{"generic", elements}
	attachment := Attachment{"template", payload}
	messages := []Message{{attachment}}
	sendMessage := SendMessage{messages}

	smbytes, _ := json.Marshal(sendMessage)
	buff := bytes.NewBuffer(smbytes)

	// Request 객체 생성
	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.11/me/message_creatives?access_token="+config.FacebookToken, buff)
	if err != nil {
		panic(err)
	}

	//Content-Type 헤더 추가
	req.Header.Add("Content-Type", "application/json")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		mci := MessageCreativeID{}
		json.Unmarshal(respBody, &mci)
		messageCreativeID = mci.ID
		return
	}
	return
}

func (noticer Noticer) sendMessage(messageCreativeID string) (broadcastID string) {
	sendBody := sendJSON{messageCreativeID, "REGULAR"}
	smbytes, _ := json.Marshal(sendBody)
	buff := bytes.NewBuffer(smbytes)

	// Request 객체 생성
	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.11/me/broadcast_messages?access_token="+config.FacebookToken, buff)
	if err != nil {
		panic(err)
	}

	//Content-Type 헤더 추가
	req.Header.Add("Content-Type", "application/json")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		bi := BroadcastID{}
		json.Unmarshal(respBody, &bi)
		broadcastID = bi.ID
		return
	}
	return
}
