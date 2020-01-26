package crawler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(b []byte) {

	d1 := b
	//[]byte("hello\ngo\n")
	err := ioutil.WriteFile("./dat1", d1, 0644)
	check(err)

	f, err := os.Create("./dat2")
	check(err)

	defer f.Close()

	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)

	w.Flush()

}

func TestMedium() {
	_ = write

	c := colly.NewCollector()

	// Find and visit all links
	/*c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	        //fmt.Println(e)
	        fmt.Println(e.Attr("href"))
			e.Request.Visit(e.Attr("href"))
		})*/

	c.OnResponse(func(r *colly.Response) {
		//fmt.Printf("%q",r.Body)
		bodyReader := bytes.NewReader(r.Body)
		//fmt.Printf("%T\n", bodyReader)
		z := html.NewTokenizer(bodyReader)
		var title string
		var finalBody string
		bodySl := make([]string, 0)

		inBody := false
		nextBody := false
		nextTitle := false

		badElem := map[string]bool{
			"script" : true,
			"noscript" : true,
			"span" : true,
			"button" : true,
			"figcaption" : true,
			"figure" : true,
			"iframe" : true,
		}

		readLoop:
		for {
			tokenType := z.Next()

			switch tokenType {
			case html.ErrorToken:
				err := z.Err()
				if err == io.EOF {
					//end of the file, break out of the loop
					break readLoop
				}
			case html.TextToken:
				tk := z.Token()
				//fmt.Println(tk.Data)
				if nextBody && strings.TrimSpace(tk.Data) != ""{
					//fmt.Println(tk.Type)
					//fmt.Println(tk.Data)
					bodySl = append(bodySl, tk.Data)
				}
				
				if nextTitle {
					//report the page title and break out of the loop
					title = tk.Data
					fmt.Println("title=", title)
					nextTitle = false
				}

			case html.StartTagToken:
				//get the token
				token := z.Token()
				dt := token.Data
				//fmt.Println(token.Data)
				if inBody && dt == "p" || dt=="a" || dt=="li" || dt=="h1" {
					//nextBody = true
				}

				if _, bad := badElem[dt]; inBody && title != "" && !bad {
					nextBody = true
				} else {
					nextBody = false
				}

				//if the name of the element is "title"
				if "title" == dt {
					//the next token should be the page title
					nextTitle = true
				}
				if dt == "body" {
					inBody = true
				}


			case html.EndTagToken:
				token := z.Token()
				//fmt.Println("/", token.Data)
				switch token.Data {
				case "article":
					break readLoop
				case "li":
					bodySl = append(bodySl, "\n")
				case "p":
					bodySl = append(bodySl, "\n")
				}
			}
		}
		spl := strings.Split(strings.Join(bodySl, " "), "\n")
		for i, s := range spl {
			spl[i] = strings.TrimSpace(s)
		}
		finalBody = strings.Join(spl, "\n")
		fmt.Println(finalBody)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://medium.com/free-code-camp/inside-the-invisible-war-for-the-open-internet-dd31a29a3f08")
	//c.Visit("http://go-colly.org/")
	//    "http://medium.com/")
}
