package crawler

import (
  "encoding/csv"
  "log"
  "os"
  "strings"
  "github.com/gocolly/colly"
)

type USToday struct {
  News
}

func (ut *USToday) Run() {
  fName := "ustoday.csv"
  file, err := os.Create(fName)
  if err != nil {
    log.Fatalf("Cannot create file %q: %s\n", fName, err)
    return
  }
  defer file.Close()
  writer := csv.NewWriter(file)
  defer writer.Flush()

  writer.Write([]string{"Title", "Body", "Time", "Url"})

  c := colly.NewCollector()

  c.OnHTML("a.gnt_m_flm_a", func(e *colly.HTMLElement) {
    title := strings.Trim(e.Text, " ")
    url := e.Attr("href")
    body := strings.Trim(e.Attr("data-c-br"), " ")
    time := e.ChildAttr("div", "data-c-dt")
    writer.Write([]string{
      title,
      body,
      time,
      url,
    })
    ut.News.title = title
    ut.News.body = body
    ut.News.time = time
    ut.News.url = url
  })

  c.Visit("https://www.usatoday.com/money/")

  log.Printf("Scraping finished, check file %q for results\n", fName)
}
