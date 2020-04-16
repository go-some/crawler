package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-some/txtanalyzer"
	"github.com/gocolly/colly"
)

type MarketWatch struct {
}

func (rc *MarketWatch) Run(wtr DocsWriter) {
	rootCollector := colly.NewCollector(
		colly.MaxDepth(2),
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
	imgSrc := ""

	articleCollector.OnHTML("head", func(e *colly.HTMLElement) {
		// cnbc의 경우 head meta 태그에 대표 이미지 정보가 저장되어 있음
		url = e.Request.URL.String()
		imgSrc = e.ChildAttr("meta[property=\"og:image\"]", "content")
	})

	articleCollector.OnHTML(".region--primary", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB

		- 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		- 크롤과 동시에 바로 저장하도록 함
		- mongoDB에서의 중복체크는 WriteDocs 함수에서 진행
		*/
		date := DateParser(e.ChildText(".timestamp "))
		// 해당 기사의 head로부터 대표 이미지를 찾고 그래프 이미지인지 check
		var hasGraphImg bool
		if url != e.Request.URL.String() || strings.Contains(imgSrc, "mw_logo_social.png") {
			imgSrc = ""
			hasGraphImg = false
		} else {
			hasGraphImg := CheckGraphImage(imgSrc)
			if !hasGraphImg {
				imgSrc = ""
			}
		}
		title := e.ChildText("h1.article__headline")
		body := e.ChildText("div.article__body ")
		entitiesInTitle, personList, orgList, prodList := txtanalyzer.NEROnDoc(title, body)
		bodySum := txtanalyzer.SumOnDoc(title, body)
		doc := News{
			Title:           title,
			Body:            body,
			Time:            date,
			Url:             e.Request.URL.String(),
			Origin:          "marketwatch",
			ImgUrl:          imgSrc,
			HasGraphImg:     hasGraphImg,
			EntitiesInTitle: entitiesInTitle,
			PersonList:      personList,
			OrgList:         orgList,
			ProdList:        prodList,
			BodySum:         bodySum,
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
