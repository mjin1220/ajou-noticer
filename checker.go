package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Notice struct {
	number     int
	title      string
	department string
	regi_date  string
}
type Checker struct {
}

func (checker Checker) check() {
	doc, err := goquery.NewDocument("http://www.ajou.ac.kr/new/ajou/notice.jsp")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tbody").Each(func(_ int, s *goquery.Selection) {
		s.Find("tr").Each(func(_ int, s *goquery.Selection) {
			notice := Notice{}
			s.Find("td").Each(func(i int, td *goquery.Selection) {
				switch {
				case i == 0:
					notice.number, err = strconv.Atoi(td.Text())
				case i == 2:
					notice.title = strings.Trim(td.Find("a").Text(), " \n	")
					// notice.title = td.Nodes
				case i == 3:
					notice.department = td.Text()
				case i == 4:
					notice.regi_date = td.Text()
				}

				//notice := Notice{}

				//fmt.Printf("Notice %d: %s - %s\n", i, band, title)
			})
			fmt.Println(notice)
			fmt.Println("-------------------------------------------------")
		})

	})
}
