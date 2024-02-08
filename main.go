package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/Machiel/slugify"
	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
)

const (
	resultsURL  = "http://www.internetculturale.it/it/16/search?instance=magindice"
	downloadURL = "http://www.internetculturale.it/metaindiceServices/MagExport?id="
	output      = "./data"
)

func pages(url string) (float64, error) {
	var body string
	err := requests.URL(url).ToString(&body).Fetch(context.Background())
	if err != nil {
		return 0, err
	}

	// pagination match, example string: `Pagina 1 di 14.671 (293.410 risultati trovati)`
	regexp, _ := regexp.Compile(`Pagina (\d+) di (\d+) \((\d+(\.\d+)?) risultati trovati\)`)
	matches := regexp.FindStringSubmatch(body)
	if len(matches) >= 3 {
		total, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return 0, err
		}
		if total == 0 {
			return 0, errors.New("0 results")
		}

		pages, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, err
		}
		return pages, nil
	} else {
		return 0, errors.New("0 results")
	}
}

func download(oai string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%s%s", downloadURL, oai)
	slug := slugify.Slugify(oai)
	filename := fmt.Sprintf("%s/%s.xml", output, slug)

	err := requests.URL(url).ToFile(filename).Fetch(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var wg sync.WaitGroup
	var q string

	query := flag.String("query", "", "query string")
	queryAll := flag.Bool("all", false, "search all (*)")
	biblioType := flag.String("biblio-type", "", "filtery by bibliographic type (eg. 'periodico')")
	documentType := flag.String("document-type", "", "filtery by document type (eg. 'manoscritto')")

	flag.Parse()
	if !*queryAll && *query == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		os.Mkdir(output, 0755)
	}

	if *queryAll {
		q = "*"
	} else {
		q = url.QueryEscape(*query)
	}

	seed := fmt.Sprintf("%s&q=%s", resultsURL, q)

	if *biblioType != "" {
		seed = seed + "&__meta_typeLivello=" + *biblioType
	}

	if *documentType != "" {
		seed = seed + "&channel__typeTipo=" + url.QueryEscape(*documentType)
	}

	logger.Info(seed)
	pages, err := pages(seed)

	if err != nil {
		logger.Warn("0 results")
		os.Exit(0)
	}

	for i := 1; i <= int(pages); i++ {
		url := fmt.Sprintf("%s&pag=%d", seed, i)

		var buf bytes.Buffer
		err = requests.URL(url).ToBytesBuffer(&buf).Fetch(context.Background())
		if err != nil {
			logger.Error(err.Error())
		}

		doc, err := goquery.NewDocumentFromReader(&buf)
		if err != nil {
			logger.Error(err.Error())
		}

		doc.Find(".dc_id").Each(func(i int, s *goquery.Selection) {
			oai := s.Text()
			slug := slugify.Slugify(oai)
			logger.Info("download", "identifier", oai, "file", slug)
			wg.Add(1)
			go download(oai, &wg)
		})

	}

	wg.Wait()

}
