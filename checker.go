package main

import (
	"fmt"
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

// func (checker Checker) makeMessage(notice Notice) {
// 	file, e := ioutil.ReadFile("./config.json")
// 	if e != nil {
// 		fmt.Printf("File error: %v\n", e)
// 		os.Exit(1)
// 	}

// 	person := Person{"Alex", 10}
// 	pbytes, _ := json.Marshal(person)
// 	buff := bytes.NewBuffer(pbytes)

// 	// Request 객체 생성
// 	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.11/me/message_creatives?access_token=", buff)
// 	if err != nil {
// 		panic(err)
// 	}

// 	//Content-Type 헤더 추가
// 	req.Header.Add("Content-Type", "application/xml")

// 	// Client객체에서 Request 실행
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// Response 체크.
// 	respBody, err := ioutil.ReadAll(resp.Body)
// 	if err == nil {
// 		str := string(respBody)
// 		println(str)
// 	}
// }
