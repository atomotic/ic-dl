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
	"strings"
	"sync"

	"github.com/Machiel/slugify"
	"github.com/PuerkitoBio/goquery"
)

const (
	resultsURL      = "http://www.internetculturale.it/it/16/search?instance=magindice"
	downloadURL     = "http://www.internetculturale.it/metaindiceServices/MagExport?id="
	outputDirectory = "./ic-data"
)

var (
	query        = flag.String("query", "", "query string")
	queryAll     = flag.Bool("all", false, "search all (*)")
	biblioType   = flag.String("biblio-type", "", "filtery by bibliographic type (eg. 'periodico')")
	documentType = flag.String("document-type", "", "filtery by document type (eg. 'manoscritto')")
)

func getPages(url string) (int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("error")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	// match the pagination, examples string `Pagina 1 di 14.671 (293.410 risultati trovati)`
	regexp, _ := regexp.Compile(`Pagina (\d+) di (\d+\.?\d*)`)

	results := regexp.FindStringSubmatch(string(body))

	if len(results) > 0 {
		// remove the dot decimal separator from number of pages, and convert to int
		pages, _ := strconv.Atoi(strings.Replace(results[2], ".", "", -1))
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
		log.Fatal("file creation error")
	}
	defer output.Close()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal("xml download error")
	}
	defer res.Body.Close()

	io.Copy(output, res.Body)
}

func main() {
	var wg sync.WaitGroup
	var startURL string
	var q string

	flag.Parse()
	if !*queryAll && *query == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		os.Mkdir(outputDirectory, 0755)
	}

	if *queryAll {
		q = "*"
	} else {
		q = url.QueryEscape(*query)
	}

	startURL = fmt.Sprintf("%s&q=%s", resultsURL, q)

	if *biblioType != "" {
		startURL = startURL + "&__meta_typeLivello=" + *biblioType
	}

	if *documentType != "" {
		startURL = startURL + "&channel__typeTipo=" + url.QueryEscape(*documentType)
	}

	fmt.Println(startURL)
	pages, err := getPages(startURL)

	if err != nil {
		fmt.Println("0 results")
		os.Exit(1)
	}

	for i := 1; i <= pages; i++ {
		url := fmt.Sprintf("%s&pag=%d", startURL, i)

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
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
