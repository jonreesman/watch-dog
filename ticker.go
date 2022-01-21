package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func (t ticker) hourlyWipe() {
	t.NumTweets = 0
	t.Tweets = nil
}

func (t ticker) computeHourlySentiment() {
	var total float64
	for _, s := range t.Tweets {
		total += float64(s.Polarity)
	}
	t.HourlySentiment = total / float64(t.NumTweets)
}

func (t ticker) pushToDb(d DBManager) {
	for _, tw := range t.Tweets {
		fmt.Println("subject:", tw.Subject)
		fmt.Println("source:", tw.Source)
		d.addStatement(tw.Expression, tw.TimeStamp, tw.Polarity)
	}
}

func (t ticker) printTicker() {
	fmt.Println("Name: ", t.Name)
	fmt.Println("Number of Tweets", t.NumTweets)
	fmt.Println("Last Scrape", t.LastScrapeTime)
	for _, tw := range t.Tweets {
		fmt.Printf("\nTimestamp: %s - Tweet: %s\n", time.Unix(tw.TimeStamp, 0).String(), tw.Expression)
		fmt.Println("Polarity: ", tw.Polarity)
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
		stock := ticker{Name: scanner.Text(), NumTweets: 0}
		if CheckTickerExists(stock.Name) {
			tick = append(tick, stock)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tick
}

func addTicker(s string) ticker {
	if !CheckTickerExists(s) {
		log.Println("Stock", s, "does not exist.")
	}
	t := ticker{
		Name: s,
	}
	t.scrape()
	return t

}

func scrapeAll(t *[]ticker) {
	for i, tick := range *t {
		(*t)[i].Tweets = append((*t)[i].Tweets, twitterScrape(tick)...)
		(*t)[i].NumTweets = len((*t)[i].Tweets)
	}
}

func (t *ticker) scrape() {
	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.NumTweets = len(t.Tweets)
}
