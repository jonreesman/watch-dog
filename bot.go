package main

import (
	"fmt"
	"time"
)

//bot.go
//fmt, time
type bot struct {
	tickers  []ticker
	interval time.Duration //defined in seconds
}

func initBot() bot {
	var b bot
	b.interval = 3600 * time.Second
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
		for i, _ := range b.tickers {
			b.tickers[i].lastScrapeTime = time.Now()
			b.tickers[i].printTicker()
			b.tickers[i].hourlySentiment = b.tickers[i].computeHourlySentiment()
			//b.tickers[i].pushToDb()	
			fmt.Printf("Sentiment: %f\n", float64(b.tickers[i].hourlySentiment))
			//b.tickers[i].dump()
		}
		time.Sleep(b.interval)
	}
}
