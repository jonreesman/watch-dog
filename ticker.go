package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func (t *ticker) hourlyWipe() {
	t.NumTweets = 0
	t.Tweets = nil
}

func (t *ticker) computeHourlySentiment() {
	var total float64
	for _, s := range t.Tweets {
		total += float64(s.Polarity)
	}
	t.HourlySentiment = total / float64(t.NumTweets)
}

func (t ticker) pushToDb(d DBManager) {
	for _, tw := range t.Tweets {
		fmt.Println("added statement to DB for:", tw.Subject)
		//fmt.Println("source:", tw.Source)
		//fmt.Println("Tweet length: ", len(tw.Expression))
		d.addStatement(tw.Expression, tw.TimeStamp, tw.Polarity, tw.PermanentURL)
	}
	d.addSentiment(t.LastScrapeTime.Unix(), t.id, t.HourlySentiment)
}

func (t ticker) printTicker() {
	fmt.Println("Name: ", t.Name)
	fmt.Println("Number of Tweets", t.NumTweets)
	fmt.Println("Last Scrape", t.LastScrapeTime)
	for _, tw := range t.Tweets {
		fmt.Printf("\nTimestamp: %s - Tweet: %s\n", time.Unix(tw.TimeStamp, 0).String(), tw.Expression)
		fmt.Println(tw.PermanentURL)
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

func addTicker(s string, d DBManager) (ticker, error) {
	if !CheckTickerExists(s) {
		log.Println("Stock", s, "does not exist.")

		return ticker{Name: "none"}, errors.New("Stock/crypto does not exist.")
	}
	t := ticker{
		Name: s,
	}
	t.id = d.addTicker(t.Name)
	t.scrape()
	t.LastScrapeTime = time.Now()
	t.computeHourlySentiment()
	t.pushToDb(d)
	return t, nil

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
