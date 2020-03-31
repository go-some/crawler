package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
)

type CNBC struct {
}

func (rc *CNBC) Run(wtr DocsWriter) {
	// Instantiate default NewCollector
	c := colly.NewCollector(
		colly.MaxDepth(1),
		// Visit business & technology section
		colly.URLFilters(
			regexp.MustCompile("https://www\\.cnbc\\.com/business"),
			regexp.MustCompile("https://www\\.cnbc\\.com/technology"),
		),
	)
	c.AllowURLRevisit = false

	// Create another collector to scrape each news article
	articleCollector := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		/* crawl all href links recursively	*/
		link := e.Request.AbsoluteURL(e.Attr("href"))
		//if the link is article page, crawl using articleCollector
		//else, visit the link until MaxDepth
		//cnbc의 기사 형식은 '/년도/일/월/제목'이므로 regxp를 활용
		re := regexp.MustCompile("https://www\\.cnbc\\.com/[0-9]{4}/[0-9]{2}/[0-9]{2}/.+")
		if re.MatchString(link) {
			articleCollector.Visit(link)
		} else {
			e.Request.Visit(link) //e.Request.Visit을 이용해야 MaxDepth 처리가 된다.
		}
	})
	// 뉴스 기사 url 별 대표 image source 를 저장하기 위한 변수 선언
	url := ""
	imgSrc := ""

	articleCollector.OnHTML("head", func(e *colly.HTMLElement) {
		// cnbc의 경우 head meta 태그에 대표 이미지 정보가 저장되어 있음
		url = e.Request.URL.String()
		imgSrc = e.ChildAttr("meta[itemprop=primaryImageOfPage]", "content")
	})

	articleCollector.OnHTML("div.Article", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB

		- 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		- 크롤과 동시에 바로 저장하도록 함
		- mongoDB에서의 중복체크는 WriteDocs 함수에서 진행
		*/
		date := DateParser(e.ChildText("time[data-testid=published-timestamp]"))
		// 해당 기사의 head로부터 대표 이미지를 잘 찾았는지 check
		if url != e.Request.URL.String() {
			imgSrc = ""
		}
		doc := News{
			Title:  e.ChildText("h1"),
			Body:   e.ChildText("div[data-module=ArticleBody]"),
			Time:   date,
			Url:    e.Request.URL.String(),
			Origin: "cnbc",
			Img:    imgSrc,
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}

	})

	c.Visit("https://www.cnbc.com/business")
}
