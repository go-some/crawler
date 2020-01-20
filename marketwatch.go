package crawler

import (
  "encoding/csv"
  "log"
  "os"
  "strings"
  "github.com/gocolly/colly"
)

type MarketWatch struct {
  News
}

func (rc *MarketWatch) Run() {
  fName := "marketwatch.csv"
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

  c.OnHTML("div.element--article", func(e *colly.HTMLElement) {
    title := strings.Trim(e.ChildText("h3.article__headline"), " ")
    url := e.ChildAttr(".article__headline a", "href")
    body := ""
    time := e.ChildText(".article__timestamp")
    writer.Write([]string{
      title,
      body,
      time,
      url,
    })
    rc.News.title = title
    rc.News.body = body
    rc.News.time = time
    rc.News.url = url
  })

  c.Visit("https://www.marketwatch.com/investing/stocks?mod=exchange-traded-funds")

  log.Printf("Scraping finished, check file %q for results\n", fName)
}
