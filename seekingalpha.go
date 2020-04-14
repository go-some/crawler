package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type SeekingAlpha struct {
}

func (rc *SeekingAlpha) Run(wtr DocsWriter) {
	// Instantiate default NewCollector
	c := colly.NewCollector(
		colly.MaxDepth(2),
		// Visit Latest News section
		colly.URLFilters(
			regexp.MustCompile("https://seekingalpha\\.com/market-outlook"),
			regexp.MustCompile("https://seekingalpha\\.com/market-news"),
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
		if strings.Index(link, "seekingalpha.com/news") != -1 || strings.Index(link, "seekingalpha.com/article") != -1 {
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
		imgSrc = e.ChildAttr("meta[property=\"og:image\"]", "content")
		if strings.Index(imgSrc, "og_image_192-59bfd51c9fe6af025b2f9f96c807e46f8e2f06c5ae787b15bf1423e6c676d4db.png") != -1 {
			imgSrc = ""
			url = ""
		}
	})

	articleCollector.OnHTML("div[id=main_content]", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB

		- 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		- 크롤과 동시에 바로 저장하도록 함
		- mongoDB에서의 중복체크는 WriteDocs 함수에서 진행
		*/
		date := DateParser(e.ChildText("time"))
		// 해당 기사의 head로부터 대표 이미지를 찾고 그래프 이미지인지 check
		var hasGraphImg bool
		if url != e.Request.URL.String() {
			imgSrc = ""
			hasGraphImg = false
		} else {
			hasGraphImg = CheckGraphImage(imgSrc)
			if !hasGraphImg {
				imgSrc = ""
			}
		}
		doc := News{
			Title:       e.ChildText("h1"),
			Body:        e.ChildText("div#a-cont"),
			Time:        date,
			Url:         e.Request.URL.String(),
			Origin:      "seekingalpha",
			ImgUrl:      imgSrc,
			HasGraphImg: hasGraphImg,
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}
	})

	c.Visit("https://seekingalpha.com/market-outlook")
}
