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
	cookie := "PHAROS_VISITOR=00002af201605039a00f6b26ca1e0013; _ga=GA1.3.1203974661.1513176736; JSESSIONID=Px3VJH1FiYF9gdrrrGSxksBR4g2Ty1hasNtOaxuQuSS3fNkkMEPj4QF2eJiVyrDQ.toegye_servlet_Portal01; ssotoken=Wruzl4OKYt9gH%2Bo9QdZbzfCPugARo68W0Jf31d%2BG7E2P7JdtwWsBQEF4AoV%2BEfeTutffqOt7APEwP7TnmD%2BdNu%2BM4lHv0fW4wsG8Loa%2FSlwSzoONlqPMRnJj7xruAta1zPC9VSFF6k1PQJLIcVpGkA%3D%3D; SSOGlobalLogouturl=get^http://portal2.ajou.ac.kr/com/sso/logout.jsp$; JSESSIONID=dkbFbmytCDJxGyHmHsdYSS0bAMxhaP8bs2X1wkOiKnHL9816M1mrFPaL7P0Y1IWK.junggak_servlet_engine2"
	req.Header.Set("Cookie", cookie)
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
