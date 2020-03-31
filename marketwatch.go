package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

type MarketWatch struct {
}

func (rc *MarketWatch) Run(wtr DocsWriter) {
	rootCollector := colly.NewCollector(
		colly.MaxDepth(1),
		colly.URLFilters(
			regexp.MustCompile("https://www\\.marketwatch\\.com/"),
			regexp.MustCompile("https://www\\.marketwatch\\.com/story/.+"),
		),
		//colly.DisallowedURLFilters(
		//	regexp.MustCompile("https://www\\.usatoday\\.com/opinion/"),
		//),
	)
	rootCollector.AllowURLRevisit = false

	articleCollector := colly.NewCollector()

	rootCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(link, "marketwatch.com/story/") != -1 {
			articleCollector.Visit(link)
		} else {
			e.Request.Visit(link)
		}
	})
	// 뉴스 기사 url 별 대표 image source 를 저장하기 위한 변수 선언
	url := ""
	img_src := ""

	articleCollector.OnHTML("head", func(e *colly.HTMLElement){
		// cnbc의 경우 head meta 태그에 대표 이미지 정보가 저장되어 있음
		url = e.Request.URL.String()
		img_src = e.ChildAttr("meta[property=\"og:image\"]", "content")
	})

	articleCollector.OnHTML(".region--primary", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB

		- 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		- 크롤과 동시에 바로 저장하도록 함
		- mongoDB에서의 중복체크는 WriteDocs 함수에서 진행
		*/
		date := dateParser(e.ChildText(".timestamp "))
		// 해당 기사의 head로부터 대표 이미지를 잘 찾았는지 check
		if url != e.Request.URL.String() || strings.Contains(img_src, "mw_logo_social.png") {
			img_src = ""
		}
		doc := News{
			Title:  e.ChildText("h1.article__headline"),
			Body:   e.ChildText("div.article__body "),
			Time:   date,
			Url:    e.Request.URL.String(),
			Origin: "MarketWatch",
			Img:	img_src,
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}

	})

	rootCollector.Visit("https://www.marketwatch.com/")
}
