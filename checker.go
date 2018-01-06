package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

const html = `
<div class="container">
    <div class="row">
      <div class="col-lg-8">
        <p align="justify"><b>Name</b>Priyaka</p>
        <p align="justify"><b>Surname</b>Patil</p>
        <p align="justify"><b>Adress</b><br>India,Kolhapur</p>
        <p align="justify"><b>Hobbies&nbsp;</b><br>Playing</p>
        <p align="justify"><b>Eduction</b><br>12th</p>
        <p align="justify"><b>School</b><br>New Highschool</p>
       </div>
    </div>
</div>
`

type Checker struct {
}

func (checker Checker) check() {
	doc, err := goquery.NewDocument("http://www.ajou.ac.kr/new/ajou/notice.jsp")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Notice %d: %s - %s\n", i, band, title)
	})
}
