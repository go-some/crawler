# crawler
crowl 서비스의 크롤러를 구현합니다.

## Preparation
- 크롤 모듈인 colly 패키지를 설치합니다 [colly](http://go-colly.org/)
```bash
go get -u github.com/gocolly/colly/...
```

## 설치
- go 설치 후
- [work-space-path]/src/ 밑에 이 프로젝트를 clone 시킵니다.
- 현재는 [work-space-path]/src/crawler 에 위치시켜놓고 쓰고있습니다.
- 만약, [work-space-path]/src/github.com/go-some/crawler 와 같이 구성하실 경우엔 executor에서 import path를 적절히 수정해주세요.
- 코드 수정 후 go install을 꼭 해주시고,
- [executor(main.go)](https://github.com/go-some/executor)를 실행해 주세요.

## 코드 설명
- 각 크롤러 코드는 해당 사이트를 의미하는 struct를 정의하고,
```go
type USToday struct {
}
```
- 특정 사이트를 크롤하는 Run 메서드를 구현 후, 리시버를 통해 struct와 매핑시킵니다.
```go
func (rc *USToday) Run(wtr Writer) {
  c := colly.NewCollector()

  docs := make([]News, 0, 100)

  c.OnHTML("a.gnt_m_flm_a", func(e *colly.HTMLElement) {
    // site-specific patterns
    title := strings.Trim(e.Text, " ")
    url := e.Attr("href")
    body := strings.Trim(e.Attr("data-c-br"), " ")
    time := e.ChildAttr("div", "data-c-dt")
    origin := "https://www.usatoday.com/"
    doc := News{ title, body, time, url, origin }
    docs = append(docs, doc)
  })

  c.OnScraped(func (r *colly.Response) {
    wtr("ustoday.csv", docs)
  })

  c.Visit("https://www.usatoday.com/money/")
}
```
- 구현된 크롤러는 [github.com/go-some/executor](https://github.com/go-some/executor)에서 호출됩니다.
