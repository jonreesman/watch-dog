package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
)

func (t ticker) computeHourlySentiment() {
	var total float64
	for _, s := range t.tweets {
		total += float64(s.polarity)
	}
	t.hourlySentiment = total / float64(t.numTweets)
}

func (t ticker) printTicker() {
	fmt.Println("Name: ", t.name)
	fmt.Println("Number of Tweets", t.numTweets)
	for _, tw := range t.tweets {
		//fmt.Println("Tweet: ", tw.expression)
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

func scrape(t *[]ticker) {
	for i, tick := range *t {
		(*t)[i].tweets = append((*t)[i].tweets, twitterScrape(tick.name)...)
		(*t)[i].numTweets = len((*t)[i].tweets)
	}
}
