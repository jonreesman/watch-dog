package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cdipaolo/sentiment"
)

func (t *ticker) hourlyWipe() {
	t.numTweets = 0
	t.Tweets = nil

}

func (t *ticker) computeHourlySentiment(sentimentModel sentiment.Models) {
	var total float64
	for _, s := range t.Tweets {
		s.Polarity = sentimentModel.SentimentAnalysis(s.Expression, sentiment.English).Score
		total += float64(s.Polarity)
	}
	t.hourlySentiment = total / float64(t.numTweets)
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

	t.hourlySentiment = total / float64(t.numTweets)
}

func (tickers *tickerSlice) pushToDb(d DBManager) {
	ts := *tickers
	for i := range ts {
		go func(t *ticker) {
			t.pushToDb(d)
			t.hourlyWipe()
		}(&ts[i])
	}
	*tickers = ts
}

func (t ticker) pushToDb(d DBManager) {
	var wg sync.WaitGroup
	for _, tw := range t.Tweets {
		fmt.Println("added statement to DB for:", tw.Subject)
		wg.Add(1)
		go d.addStatement(&wg, tw.Expression, tw.TimeStamp, tw.Polarity, tw.PermanentURL)
	}
	wg.Add(1)
	go d.addSentiment(&wg, t.lastScrapeTime.Unix(), t.Id, t.hourlySentiment)
	wg.Wait()
}

func (tickers *tickerSlice) importTickers(d DBManager) {
	existingTickers := d.returnTickers()
	ts := *tickers
	if len(existingTickers) != 0 {
		*tickers = existingTickers
		return
	}
	file, err := os.Open("tickers.txt")
	if err != nil {
		log.Panicf("faild reading data from file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tick := ticker{Name: scanner.Text(), numTweets: 0}
		if CheckTickerExists(tick.Name) {
			tick.Id, _ = d.addTicker(tick.Name)
			ts.appendTicker(tick)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	*tickers = ts
}

func (tickers *tickerSlice) addTicker(name string, d DBManager) (ticker, error) {
	s := sanitize(name)
	ts := *tickers
	if !CheckTickerExists(s) {
		log.Println("Stock does not exist.")

		return ticker{Name: "", Id: 0}, errors.New("stock/crypto does not exist")
	}
	t := ticker{
		Name: s,
	}
	var err error
	t.Id, err = d.addTicker(t.Name)
	if err != nil {
		return ticker{Name: "", Id: 0}, err
	}
	t.singleScrape()
	t.pushToDb(d)
	ts.appendTicker(t)
	*tickers = ts
	return t, nil

}

func (tickers *tickerSlice) appendTicker(t ticker) {
	ts := *tickers
	ts = append(ts, t)
	*tickers = ts
}

func (tickers *tickerSlice) deleteTicker(id int) {
	ts := *tickers
	for i := range ts {
		if ts[i].Id == id {
			ts[i] = ts[len(ts)-1]
			ts = ts[:len(ts)-1]
			break
		}
	}
	*tickers = ts

}

func (tickers *tickerSlice) scrape(sentimentModel sentiment.Models) {
	var wg sync.WaitGroup
	ts := *tickers
	for i := range ts {
		wg.Add(1)
		go ts[i].scrape(&wg, sentimentModel)
	}
	wg.Wait()
	*tickers = ts
}

func (t *ticker) scrape(wg *sync.WaitGroup, sentimentModel sentiment.Models) {
	defer wg.Done()

	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.numTweets = len(t.Tweets)
	t.lastScrapeTime = time.Now()
	t.computeHourlySentiment(sentimentModel)
}

func (t *ticker) singleScrape() {
	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.numTweets = len(t.Tweets)
	t.lastScrapeTime = time.Now()
	t.singleComputeHourlySentiment()
}
