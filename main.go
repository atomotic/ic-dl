package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/Machiel/slugify"
	"github.com/PuerkitoBio/goquery"
)

const (
	resultsURL      = "http://www.internetculturale.it/it/16/search?instance=magindice&q="
	downloadURL     = "http://www.internetculturale.it/metaindiceServices/MagExport?id="
	outputDirectory = "./ic-data"
)

var (
	query = flag.String("query", "", "string to search in Internet Culturale")
)

func getPages(query string) (int, error) {
	url := fmt.Sprintf("%s%s", resultsURL, url.QueryEscape(query))

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	regexp, _ := regexp.Compile(`Pagina (\d+) di (\d+)`)
	results := regexp.FindStringSubmatch(string(body))

	if len(results) > 0 {
		pages, _ := strconv.Atoi(results[2])
		return pages, nil
	} else {
		return 0, errors.New("no results")
	}
}

func downloadXML(oai string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%s%s", downloadURL, oai)
	slug := slugify.Slugify(oai)

	filename := fmt.Sprintf("%s/%s.xml", outputDirectory, slug)
	output, err := os.Create(filename)
	if err != nil {
		log.Fatal("error creating")
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("error get")
	}
	defer response.Body.Close()

	io.Copy(output, response.Body)
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()
	if *query == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		os.Mkdir(outputDirectory, 0755)
	}

	pages, err := getPages(*query)

	if err != nil {
		fmt.Println("no result")
		os.Exit(1)
	}

	for i := 1; i <= pages; i++ {
		url := fmt.Sprintf("%s%s&pag=%d", resultsURL, url.QueryEscape(*query), i)

		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".dc_id").Each(func(i int, s *goquery.Selection) {
			oai := s.Text()
			slug := slugify.Slugify(oai)

			fmt.Printf("%s\t%s.xml\n", oai, slug)
			wg.Add(1)
			go downloadXML(oai, &wg)
		})

	}

	wg.Wait()

}
