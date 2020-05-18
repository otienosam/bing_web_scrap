package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/otienosam/bing_web_scrap/metadata"
	"github.com/PuerkitoBio/goquery"
)

func handler(i int, s *goquery.Selection) {
	url, ok := s.Find("a").Attr("href")// extracts the href from the
  // goquery instance
	if !ok {
		return
	}

	fmt.Printf("%d: %s\n", i, url)
	res, err := http.Get(url)//retrives the document from the url
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(res.Body)// reads the response body
	if err != nil {
		return
	}
	defer res.Body.Close()

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))//creates
  // a zip reader
	if err != nil {
		return
	}

	cp, ap, err := metadata.NewProperties(r)// the function consumes a zip reader
	if err != nil {
		return
	}

	log.Printf(
		"%21s %s - %s %s\n",
		cp.Creator,
		cp.LastModifiedBy,
		ap.Application,
		ap.GetMajorVersion())
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Missing required argument. Usage: main.go <domain> <ext>")
	}
	domain := os.Args[1]
	filetype := os.Args[2]

	q := fmt.Sprintf(// filter string
		"site:%s && filetype:%s && instreamset:(url title):%s",
		domain,
		filetype,
		filetype)
  // filter string encoded and builds the search url
	search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))
	doc, err := goquery.NewDocument(search)// sends the search query implicitly
	if err != nil {
		log.Panicln(err)
	}
  // element selector to iterate over the document and finds a match
	s := "html body div#b_content ol#b_results li.b_algo div.b_title h2"
  // goquery inspects the document
	doc.Find(s).Each(handler)// for each matching element a handle function is called
}
