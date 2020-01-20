
package crawler

import (
	"encoding/csv"
	"log"
	"os"
	"github.com/gocolly/colly"
)

type CryptoCointMarket struct {
}

func (ccm CryptoCointMarket) Run() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".cmc-table-row", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("td.cmc-table__cell--sort-by__name"),
			e.ChildText("td.cmc-table__cell--sort-by__symbol"),
			e.ChildText("td.cmc-table__cell--sort-by__price"),
			e.ChildText("td.cmc-table__cell--sort-by__volume-24-h"),
			e.ChildText("td.cmc-table__cell--sort-by__market-cap"),
		})

	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
