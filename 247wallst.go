package crawler

import (
	"fmt"
	"regexp"

	"github.com/go-some/txtanalyzer"
	"github.com/gocolly/colly"
	nats "github.com/nats-io/nats.go"
)

type WallST247 struct {
	webCollector     *colly.Collector
	articleCollector *colly.Collector
}

func (rc *WallST247) Init() {
	rc.webCollector = colly.NewCollector(
		colly.MaxDepth(2),
		colly.URLFilters(
			regexp.MustCompile("https://247wallst\\.com/"),
		),
	)
	rc.webCollector.AllowURLRevisit = false

	rc.articleCollector = colly.NewCollector()
}

func (rc *WallST247) WebSurfing(nc *nats.Conn) {
	rc.webCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		/* crawl all href links recursively	*/
		link := e.Request.AbsoluteURL(e.Attr("href"))
		//if the link is article page, crawl using articleCollector
		//else, visit the link until MaxDepth
		re := regexp.MustCompile("https://247wallst\\.com/[a-z-]+/[0-9]{4}/[0-9]{2}/[0-9]{2}/.+")
		if re.MatchString(link) {
			//err := wtr.CheckDuplicate(link)
			//if err == nil {
			//	fmt.Printf("Already exist (%s)\n", link)
			//} else {
			//	fmt.Println(link)
			//}
			fmt.Println(link)
			nc.Publish("new_url", []byte(link))
			nc.Flush()
		} else {
			e.Request.Visit(link) //e.Request.Visit을 이용해야 MaxDepth 처리가 된다.
		}
	})
	rc.webCollector.Visit("https://247wallst.com/")
}

func (rc *WallST247) Crowl(wtr DocsWriter) {
	//Queue 에서 읽어서 크롤하는 부분
}

func (rc *WallST247) Run(wtr DocsWriter) {
	// Instantiate default NewCollector
	c := colly.NewCollector(
		colly.MaxDepth(2),
		// Crawl from the main page
		colly.URLFilters(
			regexp.MustCompile("https://247wallst\\.com/"),
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
		//cnbc의 기사 형식은 '카테고리/년도/일/월/제목'이므로 regxp를 활용
		re := regexp.MustCompile("https://247wallst\\.com/[a-z-]+/[0-9]{4}/[0-9]{2}/[0-9]{2}/.+")
		if re.MatchString(link) {
			err := wtr.CheckDuplicate(link)
			if err == nil {
				fmt.Printf("Already exist (%s)\n", link)
			} else {
				articleCollector.Visit(link)
			}

		} else {
			e.Request.Visit(link) //e.Request.Visit을 이용해야 MaxDepth 처리가 된다.
		}
	})

	articleCollector.OnHTML("div.primary", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB

		- 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		- 크롤과 동시에 바로 저장하도록 함
		- mongoDB에서의 중복체크는 WriteDocs 함수에서 진행
		*/
		date := DateParser(e.ChildText("div.post-date"))
		title := e.ChildText("div.title")
		body := e.ChildText("p")
		entitiesInTitle, personList, orgList, prodList := txtanalyzer.NEROnDoc(title, body)
		bodySum := txtanalyzer.SumOnDoc(title, body)
		doc := News{
			Title:           title,
			Body:            body,
			Time:            date,
			Url:             e.Request.URL.String(),
			Origin:          "247wallst",
			ImgUrl:          "",
			HasGraphImg:     false,
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

	c.Visit("https://247wallst.com/")
}
