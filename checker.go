package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Checker 타입은 공지사항 슬라이스들을 가짐
type Checker struct {
	OldNotices Notices
	NewNotices Notices
}

// Notices 타입은 Notice 타입의 슬라이스
type Notices []Notice

// Notice 타입은 공지사항의 항목들을 가지는 구조체 변수
type Notice struct {
	Number     int
	Title      string
	URL        string
	Department string
	RegiDate   string
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

	q := `#jwxe_main_content > div > div.list_wrap > table > tbody > tr`
	doc.Find(q).Each(func(i int, s *goquery.Selection) {
		tempNotice := Notice{}

		tempNotice.Number, _ = strconv.Atoi(strings.TrimSpace(s.Find(`td:nth-child(1)`).Text()))
		tempNotice.Title = strings.TrimSpace(s.Find(`td:nth-child(3) > a`).Text())
		tempNotice.URL, _ = s.Find(`td:nth-child(3) > a`).Attr("href")
		tempNotice.URL = url + tempNotice.URL
		tempNotice.Department = strings.TrimSpace(s.Find(`td:nth-child(4)`).Text())
		tempNotice.RegiDate = strings.TrimSpace(s.Find(`td:nth-child(5)`).Text())

		checker.NewNotices = append(checker.NewNotices, tempNotice)
	})

	// doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
	// 	s.Find("tr").Each(func(index int, s *goquery.Selection) {
	// 		tempNotice := Notice{}
	// 		s.Find("td").Each(func(i int, td *goquery.Selection) {
	// 			switch {
	// 			case i == 0:
	// 				tempNotice.Number, err = strconv.Atoi(td.Text())
	// 			case i == 2:
	// 				tempNotice.Title = strings.Trim(td.Find("a").Text(), " \n	")
	// 				tempNotice.URL, _ = td.Find("a").Attr("href")
	// 				tempNotice.URL = url + tempNotice.URL
	// 			case i == 3:
	// 				tempNotice.Department = td.Text()
	// 			case i == 4:
	// 				tempNotice.RegiDate = td.Text()
	// 			}
	// 		})
	// 		checker.NewNotices = append(checker.NewNotices, tempNotice)
	// 	})
	// })

	// if checker.OldNotices == nil { // in first check
	// 	checker.OldNotices = checker.NewNotices
	// 	fmt.Println("[ajou-noticer] first start")
	// 	return
	// }

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
