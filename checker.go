package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Notice struct {
	number     int
	title      string
	url        string
	department string
	regi_date  string
}
type Checker struct {
	notice Notice
}

type SendMessage struct {
	Messages []Message `json:"messages"`
}
type Message struct {
	Attachment Attachment `json:"attachment"`
}
type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}
type Payload struct {
	TemplateType string    `json:"template_type"`
	Elements     []Element `json:"elements"`
}
type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Buttons  []Button `json:"buttons"`
}
type Button struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type MessageCreativeId struct {
	ID string `json:"message_creative_id"`
}

type sendJSON struct {
	MessageCreativeId string `json:"message_creative_id"`
	NotificationType  string `json:"notification_type"`
}

type BroadcastId struct {
	ID string `json:"broadcast_id"`
}

func (checker Checker) check() {
	url := "http://www.ajou.ac.kr/new/ajou/notice.jsp"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
	}
	req.Header.Set("Cookie", config.PortalCookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.ajou.ac.kr")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()
	root, err := html.Parse(res.Body)
	if err != nil {
		return
	}
	doc := goquery.NewDocumentFromNode(root)

	doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
		s.Find("tr").Each(func(_ int, s *goquery.Selection) {
			s.Find("td").Each(func(i int, td *goquery.Selection) {
				switch {
				case i == 0:
					checker.notice.number, err = strconv.Atoi(td.Text())
				case i == 2:
					checker.notice.title = strings.Trim(td.Find("a").Text(), " \n	")
					checker.notice.url, _ = td.Find("a").Attr("href")
					checker.notice.url = url + checker.notice.url
					// notice.title = td.Nodes
				case i == 3:
					checker.notice.department = td.Text()
				case i == 4:
					checker.notice.regi_date = td.Text()
				}
			})
			fmt.Println(checker.notice)
			fmt.Println("-------------------------------------------------")
		})

	})
}

func (checker Checker) makeMessage() (message_creative_id string) {
	buttons := []Button{{"web_url", "http://www.ajou.ac.kr/new/ajou/notice.jsp?mode=view&article_no=171773&board_wrapper=%2Fnew%2Fajou%2Fnotice.jsp&pager.offset=0&board_no=33", "자세히 보기"}}
	elements := []Element{{"AJOU-CSR 동계 통계워크샵 및 논문특강(논문특강 강의실 변경)", "사회과학대학교학팀", buttons}}
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
		mci := MessageCreativeId{}
		json.Unmarshal(respBody, &mci)
		message_creative_id = mci.ID
		fmt.Println("\nCreate Message Success!!\n", string(respBody), "\n")
		return
	}
	return
}

func (checker Checker) sendMessage(message_creative_id string) (broadcast_id string) {
	sendJson := sendJSON{message_creative_id, "REGULAR"}
	smbytes, _ := json.Marshal(sendJson)
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
		bi := BroadcastId{}
		json.Unmarshal(respBody, &bi)
		broadcast_id = bi.ID
		fmt.Println("\nSend Message Success!!\n", string(respBody), "\n")
		return
	}
	return
}
