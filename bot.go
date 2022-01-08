package main

import (
	"fmt"
	"time"
)

func initBot() bot {
	var b bot
	b.interval = 10 * time.Second
	b.tickers = importTickers()
	return b
}

func (b bot) run() {
	for _, stock := range b.tickers {
		fmt.Println(stock)
	}
	for {
		fmt.Println("Loop")
		scrapeAll(&b.tickers)
		for i := range b.tickers {
			b.tickers[i].LastScrapeTime = time.Now()
			b.tickers[i].printTicker()
			b.tickers[i].computeHourlySentiment()
			//b.tickers[i].pushToDb()
			fmt.Printf("Sentiment: %f\n", float64(b.tickers[i].HourlySentiment))
			b.tickers[i].dump_raw()
			b.tickers[i].dump_text()
			b.tickers[i].hourlyWipe()
		}
		time.Sleep(b.interval)
	}
}
