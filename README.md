# crawler
crowl 서비스의 크롤러를 구현합니다.

## Requirements
- 크롤 모듈인 colly 패키지를 설치합니다 [colly](http://go-colly.org/)
```bash
go get -u github.com/gocolly/colly/...
```
- 몽고디비를 위한 라이브러리를 설치합니다
```bash
go get go.mongodb.org/mongo-driver
```
- [mongodb-with-go-tutorial](https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial)
- DNS error 발생시 /etc/resolve.conf 수정 [참조](https://stackoverflow.com/questions/55660134/cant-connect-to-mongo-cloud-mongodb-database-in-golang-on-ubuntu)
- [TODO] dep라고 하는 패키지 매니저가 있다고 합니다... 한 번 알아봐야 할 듯!

## Installation 
- go 설치 후
- [workspace-path] 밑에서 go-some crawler를 설치합니다
```bash
go get -u github.com/go-some/crawler
```
- .bashrc에 DBID, DBPW, DBADDR 환경변수를 정의합니다. (따로 문의)
- 코드 수정 후 go install을 꼭 해주시고,
- [executor(main.go)](https://github.com/go-some/executor)를 실행해 주세요.

## Structure
- 각 크롤러 코드는 해당 사이트를 의미하는 struct를 정의하고,
```go
type Reuters struct {
}
```
- 해당 사이트를 탐색하면서 기사(article)들을 크롤하는 Run 메서드를 구현합니다.
- 최종적으로 크롤해야 하는 기사들의 url 형식을 정의하고 `articleCollector`의 리시버를 통해 저장합니다.
- 기사 내용은 `News` struct 형식에 맞게 mongoDB에 저장되며 `WriteDocs`(writer.go)함수에서 그 기능을 수행합니다.
```go
func (rc *Reuters) Run(wtr DocsWriter) {
	// Instantiate default NewCollector
	c := colly.NewCollector(
		colly.MaxDepth(3),
		// Visit only finance and businessnews section
		colly.URLFilters(
			regexp.MustCompile("https://www\\.reuters\\.com/finance"),
			regexp.MustCompile("https://www\\.reuters\\.com/news/archive/businessnews.+"),
		),
		colly.DisallowedURLFilters(
			regexp.MustCompile("https://www\\.reuters\\.com/finance/.+"),
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
		if strings.Index(link, "reuters.com/article") != -1 {
			articleCollector.Visit(link)
		} else {
			e.Request.Visit(link) //e.Request.Visit을 이용해야 MaxDepth 처리가 된다.
		}
	})

	articleCollector.OnHTML("div.StandardArticle_inner-container", func(e *colly.HTMLElement) {
		/* Read article page and save to mongoDB
    
		  - 최종적으로 우리가 크롤하고자 하는 기사 페이지 (leaf node)
		*/
		doc := News{
			Title:  e.ChildText(".ArticleHeader_headline"),
			Body:   e.ChildText("div.StandardArticleBody_body"),
			Time:   e.ChildText(".ArticleHeader_date"),
			Url:    e.Request.URL.String(),
			Origin: "Reuters",
		}
		cnt, err := wtr.WriteDocs([]News{doc})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(cnt, "docs saved")
		}
	})

	c.Visit("https://www.reuters.com/finance")
}
```
- 구현된 크롤러는 [github.com/go-some/executor](https://github.com/go-some/executor)에서 호출됩니다.

## News Company List
- Reuters
- USAToday
- SeekingAlpha
- CNBC
- 24/7WallST
