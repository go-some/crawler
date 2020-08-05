package main

import (
	"fmt"

	"github.com/go-some/crawler"
	nats "github.com/nats-io/nats.go"
)

type Crawler interface {
	Init()
	//Run(crawler.DocsWriter)
	WebSurfing(*nats.Conn)
}

func main() {
	// 각 사이트의 크롤러를 등록
	crawlers := []Crawler{
		//&crawler.MarketWatch{},
		//&crawler.Reuters{},
		//&crawler.SeekingAlpha{},
		//&crawler.CNBC{},
		&crawler.WallST247{},
		//&crawler.USAToday{},
		/* 여기에 추가 해주세요*/
	}
	// nats server
	const natsURL = "nats://127.0.0.1:4171"
	nc, _ := nats.Connect(natsURL)

	fmt.Println("Run Crawler")

	// mongoDB writer의 구현체를 얻음
	// crawler 패키지의 writer.go에 interface DocsWriter를 구현하는 구현체들이 모여 있음

	//wtr := crawler.NewMongoDBWriter()

	// connection pool (collection handler) 생성
	//err := wtr.Init()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 크롤러의 실제 구현을 이용해 실행시키는 부분
	for _, crawler := range crawlers {
		crawler.Init()
		crawler.WebSurfing(nc)
	}

	fmt.Println("Fin Crawler")
}
