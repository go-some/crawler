# crawler
crowl 서비스의 크롤러를 구현합니다.

- 각 크롤러 코드는 해당 사이트를 의미하는 struct를 정의하고,
```go
type USToday struct {
  News // 크롤된 데이터 관리를 위해서 임시로 넣어봤습니다 (go에서의 상속, struct embedding).
}
```
- 특정 사이트를 크롤하는 Run 메서드를 구현 후, 리시버를 통해 매핑시킵니다.
```go
func (ut *USToday) Run() {
  // 생략...
  c := colly.NewCollector()

  c.OnHTML("a.gnt_m_flm_a", func(e *colly.HTMLElement) {
    title := strings.Trim(e.Text, " ")
    url := e.Attr("href")
    body := strings.Trim(e.Attr("data-c-br"), " ")
    time := e.ChildAttr("div", "data-c-dt")
    // 생략 ...
    ut.News.title = title
    ut.News.body = body
    ut.News.time = time
    ut.News.url = url
  })

  c.Visit("https://www.usatoday.com/money/")
}
```
- 구현된 크롤러는 github.com/go-some/executor에서 호출됩니다.
