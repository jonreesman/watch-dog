package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (t ticker) dump_raw() {
	p, err := json.MarshalIndent(t, "", "	")
	if err != nil {
		log.Panicf("Error dumping %s\n", err)
	}
	directory := "json/"
	filename := t.Name + "-" + strconv.FormatInt(t.lastScrapeTime.Unix(), 10)
	err = os.WriteFile(directory+filename+".json", p, 0644)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
}

func (t ticker) dump_text() {
	directory := "raw_text/"
	filename := t.Name + ".txt"
	p, err := os.OpenFile(directory+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicf("Error dumping %s\n", err)
	}
	defer p.Close()
	for _, tweet := range t.Tweets {
		printedExpression := strings.Replace(tweet.Expression, "\n", "|||", -1)
		if _, err := p.WriteString(printedExpression + "\n"); err != nil {
			log.Panicf("failed reading data from file: %s", err)
		}
	}
}

func (t ticker) printTicker() {
	fmt.Println("Name: ", t.Name)
	fmt.Println("Last Scrape", t.lastScrapeTime)
	for _, tw := range t.Tweets {
		fmt.Printf("\nTimestamp: %s - Tweet: %s\n", time.Unix(tw.TimeStamp, 0).String(), tw.Expression)
		fmt.Println(tw.PermanentURL)
		fmt.Println("Polarity: ", tw.Polarity)
	}
}
