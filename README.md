# crawler
crowl 서비스의 크롤러를 구현합니다.

## Preparation
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
type USToday struct {
}
```
- 특정 사이트를 크롤하는 Run 메서드를 구현 후, 리시버를 통해 struct와 매핑시킵니다.
- Run 함수에서는 대상 문서 리스트를 크롤하고 DocsWriter를 통해 결과를 기록합니다.
```go
func (rc *USToday) Run(wtr DocsWriter) {
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
    wtr.WriteDocs(docs)
  })

  c.Visit("https://www.usatoday.com/money/")
}
```
- 구현된 크롤러는 [github.com/go-some/executor](https://github.com/go-some/executor)에서 호출됩니다.
