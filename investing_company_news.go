
package main

import (

	"fmt"
	"github.com/gocolly/colly"
)

func main() {

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".mediumTitle1", func(e *colly.HTMLElement) {
    var titles []string

    // There are non-related news which have 'ga-label' label names "Popular News - Article"
    if e.Attr("ga-label") != "Popular News - Article"{
      titles = e.ChildTexts("a.title")
    }
    // Print new titles
    for _, title  := range titles {
  		fmt.Println(title)
  	}

	})

	c.Visit("https://www.investing.com/equities/apple-computer-inc-news")

}
