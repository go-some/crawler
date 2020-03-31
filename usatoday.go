package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

type USAToday struct {
}

func (rc *USAToday) Run(wtr DocsWriter) {
	rootCollector := colly.NewCollector(
		colly.MaxDepth(1),
		colly.URLFilters(
			regexp.MustCompile("https://www\\.usatoday\\.com/money/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/tech/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/story/money/.+"),
			regexp.MustCompile("https://www\\.usatoday\\.com/story/tech/.+"),
		),
		colly.DisallowedURLFilters(
			regexp.MustCompile("https://www\\.usatoday\\.com/news/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/sports/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/entertainment/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/news/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/life/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/travel/"),
			regexp.MustCompile("https://www\\.usatoday\\.com/opinion/"),
		),
	)
	rootCollector.AllowURLRevisit = false

	articleCollector := colly.NewCollector()

	rootCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(link, "usatoday.com/story/money") != -1 {
			articleCollector.Visit(link)
		} else {
			e.Request.Visit(link)
		}
	})
	// 뉴스 기사 url 별 대표 image source 를 저장하기 위한 변수 선언
	url := ""
	imgSrc := ""

	articleCollector.OnHTML("head", func(e *colly.HTMLElement) {
		// cnbc의 경우 head meta 태그에 대표 이미지 정보가 저장되어 있음
		url = e.Request.URL.String()
		imgSrc = e.ChildAttr("meta[property=\"og:image\"]", "content")
		idx := strings.Index(imgSrc, "?")
		imgSrc = imgSrc[:idx]
	})

	articleCollector.OnHTML("main.gnt_cw", func(e *colly.HTMLElement) {
		date := DateParser(e.ChildAttr(".gnt_ar_dt", "aria-label"))
		// 해당 기사의 head로부터 대표 이미지를 잘 찾았는지 check
		if url != e.Request.URL.String() {
			imgSrc = ""
		}
		doc := News{
			Title:  e.ChildText("h1.gnt_ar_hl"),
			Body:   e.ChildText("div.gnt_ar_b"),
			Time:   date,
			Url:    e.Request.URL.String(),
			Origin: "Usatoday",
			Img:    imgSrc,
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}
	})

	rootCollector.Visit("https://www.usatoday.com/money/")
}
