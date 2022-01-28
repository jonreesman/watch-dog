package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cdipaolo/sentiment"
)

func (t *ticker) hourlyWipe() {
	t.NumTweets = 0
	t.Tweets = nil
}

func (t *ticker) computeHourlySentiment(sentimentModel sentiment.Models) {
	var total float64
	for _, s := range t.Tweets {
		s.Polarity = sentimentModel.SentimentAnalysis(s.Expression, sentiment.English).Score
		total += float64(s.Polarity)
	}
	t.HourlySentiment = total / float64(t.NumTweets)
}

func (t *ticker) singleComputeHourlySentiment() {
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Print("Error loading sentiment analysis model")
	}
	var total float64
	for _, s := range t.Tweets {
		s.Polarity = sentimentModel.SentimentAnalysis(s.Expression, sentiment.English).Score
		total += float64(s.Polarity)
	}
	t.HourlySentiment = total / float64(t.NumTweets)
}

func (t ticker) pushToDb(d DBManager) {
	for _, tw := range t.Tweets {
		//fmt.Println("added statement to DB for:", tw.Subject)
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

func importTickers(d DBManager) []ticker {
	var tickers []ticker
	existingTickers := d.retrieveTickers()
	for i, name := range existingTickers {
		if name == "" {
			continue
		}
		tickers = append(tickers, ticker{Name: name, id: i})
		fmt.Printf("Added ticker %s with id %d", name, i)
	}
	if len(tickers) != 0 {
		return tickers
	}
	file, err := os.Open("tickers.txt")
	if err != nil {
		log.Panicf("faild reading data from file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tick := ticker{Name: scanner.Text(), NumTweets: 0}
		if CheckTickerExists(tick.Name) {
			tick.id, _ = d.addTicker(tick.Name)
			tickers = append(tickers, tick)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tickers
}

func addTicker(s string, d DBManager) (ticker, error) {
	if !CheckTickerExists(s) {
		log.Println("Stock", s, "does not exist.")

		return ticker{Name: "none"}, errors.New("stock/crypto does not exist")
	}
	t := ticker{
		Name: s,
	}
	var err error
	t.id, err = d.addTicker(t.Name)
	if err != nil {
		return ticker{Name: "none"}, err
	}
	t.scrape()
	t.LastScrapeTime = time.Now()
	t.singleComputeHourlySentiment()
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
