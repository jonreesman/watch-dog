package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cdipaolo/sentiment"
)

func (t *ticker) hourlyWipe() {
	t.numTweets = 0
	t.Tweets = nil

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
		go d.addStatement(&wg, t.Id, tw.Expression, tw.TimeStamp, tw.Polarity, tw.PermanentURL)
	}
	wg.Add(1)
	go d.addSentiment(&wg, t.LastScrapeTime.Unix(), t.Id, t.hourlySentiment)
	wg.Add(1)
	go d.updateTicker(&wg, t.Id, t.LastScrapeTime)
	wg.Wait()
}

func (tickers *tickerSlice) importTickers(d DBManager) {
	existingTickers := d.returnActiveTickers()
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
			fmt.Println("Exists")
			tick.Id, _ = d.addTicker(tick.Name)
			ts.appendTicker(tick)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	*tickers = ts
}

func (tickers *tickerSlice) appendTicker(t ticker) {
	ts := *tickers
	ts = append(ts, t)
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
	t.LastScrapeTime = time.Now()
	wg.Add(1)
	t.computeHourlySentiment(wg, sentimentModel)
}

type pack struct {
	Tweet string `json:"tweet"`
}

func (t *ticker) computeHourlySentiment(wg *sync.WaitGroup, sentimentModel sentiment.Models) {
	defer wg.Done()
	var total float64
	var response float64
	for _, s := range t.Tweets {
		p := pack{Tweet: s.Expression}
		js, _ := json.Marshal(p)
		fmt.Println(p)
		pythonSentiment, err := http.Post("http://localhost:8000/tweet", "application/json", bytes.NewBuffer(js))
		resp, err := ioutil.ReadAll(pythonSentiment.Body)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal([]byte(resp), &response)
		//s.Polarity = sentimentModel.SentimentAnalysis(s.Expression, sentiment.English).Score
		//fmt.Println("Python:", pythonSentiment.Body, " Go:", s.Polarity)

		total += float64(response)
	}
	t.hourlySentiment = total / float64(t.numTweets)
}

//DEPRECATED
func (t *ticker) singleScrape() {
	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.numTweets = len(t.Tweets)
	t.LastScrapeTime = time.Now()
	t.singleComputeHourlySentiment()
}

//DEPRECATED
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
