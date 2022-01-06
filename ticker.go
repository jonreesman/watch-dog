package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"time"
)

func (t ticker) computeHourlySentiment() float64 {
	var total float64
	for _, s := range t.tweets {
		total += float64(s.polarity)
	}
	return total / float64(t.numTweets)
}

func (t ticker) pushToDb() {
	for _, tw:= range t.tweets {
		fmt.Println("subject:",tw.subject)
		fmt.Println("source:",tw.source)
		dbPush(tw)
	}
}

func (t ticker) printTicker() {
	fmt.Println("Name: ", t.name)
	fmt.Println("Number of Tweets", t.numTweets)
	fmt.Println("Last Scrape",t.lastScrapeTime)
	for _, tw := range t.tweets {
		fmt.Printf("\nTimestamp: %s - Tweet: %s\n", time.Unix(tw.timeStamp,0).String(), tw.expression)
		fmt.Println("Polarity: ", tw.polarity)
	}
}

func importTickers() []ticker {
	file, err := os.Open("tickers.txt")
	if err != nil {
		log.Panicf("faild reading data from file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tick []ticker
	for scanner.Scan() {
		stock := ticker{name: scanner.Text(), numTweets: 0}
		tick = append(tick, stock)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tick
}

func scrapeAll(t *[]ticker) {
	for i, tick := range *t {
		(*t)[i].tweets = append((*t)[i].tweets, twitterScrape(tick.name)...)
		(*t)[i].numTweets = len((*t)[i].tweets)
	}
}

func (t ticker) scrape () {
	t.tweets = append(t.tweets, twitterScrape(t.name)...)
	t.numTweets = len(t.tweets)
}
