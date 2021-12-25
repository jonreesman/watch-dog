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
		scrape(&b.tickers)
		for i, _ := range b.tickers {
			b.tickers[i].printTicker()
			b.tickers[i].computeHourlySentiment()
			fmt.Println("Sentiment: ", "%f/n", b.tickers[i].hourlySentiment)
		}
		time.Sleep(b.interval)
	}
}
