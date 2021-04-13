package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"rfa/bq"
	"rfa/twitter"
	"rfa/vision_texts"
	"strconv"
	"time"

	"google.golang.org/api/googleapi"
)

func main() {
	projectID := flag.String("p", "", "gcp_project_id")
	location := flag.String("l", "us", "bigquery_location")
	twitterId := flag.String("u", "", "twitter_id")
	sizeStr := flag.String("s", "1", "search_size")
	flag.Parse()

	size, _ := strconv.Atoi(*sizeStr)

	var lastExecutedAt time.Time
	latest, err := bq.GetLatest(*projectID, *location, *twitterId)
	if err != nil {
		var gerr *googleapi.Error
		if ok := errors.As(err, &gerr); ok {
			switch gerr.Code {
			case 404:
				lastExecutedAt = time.Time{}
			default:
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		lastExecutedAt = latest[0].CreatedAt
	}

	rslts := twitter.Search(twitterId, size, lastExecutedAt)

	for _, rslt := range rslts {
		urls := rslt.MediaUrlHttps
		for _, url := range urls {
			fmt.Println(url)
			file := twitter.GetImage(url)
			defer os.Remove(file.Name())

			text := vision_texts.Detect(file.Name())
			if text == "" {
				continue
			}
			csvName := bq.CreateCsv(*twitterId, rslt.CreatedAt, url, text)
			defer os.Remove(csvName)
			if csvName == "" {
				continue
			}

			err := bq.LoadCsv(*projectID, csvName)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
