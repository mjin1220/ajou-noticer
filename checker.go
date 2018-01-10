package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Checker struct {
	OldNotices Notices
	NewNotices Notices
}

type Notices []Notice
type Notice struct {
	Number     int
	Title      string
	URL        string
	Department string
	RegiDate   string
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

type MessageCreativeID struct {
	ID string `json:"message_creative_id"`
}
type BroadcastID struct {
	ID string `json:"broadcast_id"`
}

type sendJSON struct {
	MessageCreativeID string `json:"message_creative_id"`
	NotificationType  string `json:"notification_type"`
}

func (checker *Checker) check() {
	url := "http://www.ajou.ac.kr/new/ajou/notice.jsp"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
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
		panic(err)
	}
	defer res.Body.Close()
	root, err := html.Parse(res.Body)
	if err != nil {
		return
	}
	doc := goquery.NewDocumentFromNode(root)

	doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
		s.Find("tr").Each(func(index int, s *goquery.Selection) {
			tempNotice := Notice{}
			s.Find("td").Each(func(i int, td *goquery.Selection) {
				switch {
				case i == 0:
					tempNotice.Number, err = strconv.Atoi(td.Text())
				case i == 2:
					tempNotice.Title = strings.Trim(td.Find("a").Text(), " \n	")
					tempNotice.URL, _ = td.Find("a").Attr("href")
					tempNotice.URL = url + tempNotice.URL
				case i == 3:
					tempNotice.Department = td.Text()
				case i == 4:
					tempNotice.RegiDate = td.Text()
				}
			})
			checker.NewNotices = append(checker.NewNotices, tempNotice)
		})
	})
	if checker.OldNotices == nil { // in first check
		checker.OldNotices = checker.NewNotices
		fmt.Println("[ajou-noticer] first start")
		return
	}

	diffNotices := checker.diff()
	if len(diffNotices) != 0 {
		new(Noticer).Notify(diffNotices)
	}
	checker.OldNotices = checker.NewNotices
	checker.NewNotices = Notices{}
	return
}

func (checker Checker) diff() (diffNotices Notices) {
	for i := 0; i < len(checker.NewNotices); i++ {
		if checker.OldNotices.contain(checker.NewNotices[i]) == false {
			diffNotices = append(diffNotices, checker.NewNotices[i])
		}
	}
	return
}

func (notices Notices) contain(notice Notice) (ret bool) {
	for i := 0; i < len(notices); i++ {
		if (notice.Number == notices[i].Number) && (notice.Title == notices[i].Title) {
			ret = true
			return
		}
	}
	ret = false
	return
}
