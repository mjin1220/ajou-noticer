package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Noticer 타입은 Notify()와 makeMessage(), sendMessage()를 가질 구조체
type Noticer struct {
}

// SendMessage 타입은 Message들을 가지고 있는 구조체
type SendMessage struct {
	Messages []Message `json:"messages"`
}

// Message 타입은 Attachment를 가지고 있는 구조체
type Message struct {
	Attachment Attachment `json:"attachment"`
}

// Attachment 타입은 Message 내 첨부할 구조체
type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

// Payload 타입은 실제로 보낼 데이터들이 있는 구조체
type Payload struct {
	TemplateType string    `json:"template_type"`
	Elements     []Element `json:"elements"`
}

// Element 타입은 하나의 요소가 담겨져있는 구조체
type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Buttons  []Button `json:"buttons"`
}

// Button 타입은 Element 내 버튼의 요소가 정의된 구조체
type Button struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// MessageCreativeID 타입은 메세지 생성요청의 응답으로 반환하는 JSON을 저장할 구조체
type MessageCreativeID struct {
	ID string `json:"message_creative_id"`
}

// BroadcastID 타입은 메세지 전송요청의 응답으로 반환하는 JSON을 저장할 구조체
type BroadcastID struct {
	ID string `json:"broadcast_id"`
}

// sendJSON 타입은 메세지 전송요청에 담아보낼 JSON을 저장할 구조체
type sendJSON struct {
	MessageCreativeID string `json:"message_creative_id"`
	NotificationType  string `json:"notification_type"`
}

// Notify 함수는 매개변수로 받은 Notices를 Facebook Message로 모두 보내는 함수
func (noticer Noticer) notify(notices Notices) {
	for i := len(notices) - 1; i >= 0; i-- {
		noticer.sendMessage(noticer.makeMessage(notices[i]))
		fmt.Printf("[%s] %s\n", "New", notices[i].Title)
	}
}

// makeMessage 함수는 Facebook message를 생성하는 함수
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

// sendMessage 함수는 Facebook message를 보내는 함수
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
