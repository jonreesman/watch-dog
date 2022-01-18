package main

import (
	"fmt"
	"time"

	gcloud "watch-dog/databases/gclouddriver"
)

func initBot() bot {
	var b bot
	b.interval = 3600 * time.Second
	b.tickers = importTickers()
	return b
}

func (b bot) run() {
	var d gcloud.DBManager
	d.InitializeManager()
	if d.db == nil {
		fmt.Println("Uh oh")
	}
	d.dropTable("tickers")
	d.dropTable("statements")
	fmt.Println("Creating ticker table")
	d.createTickerTable()
	fmt.Println("Creating stmt table")
	d.createStatementTable()
	fmt.Println("Statement created")
	for _, stock := range b.tickers {
		d.addTicker(stock)
		fmt.Println(stock)
	}
	d.retrieveTickers()
	for {
		fmt.Println("Loop")
		scrapeAll(&b.tickers)
		for i := range b.tickers {
			b.tickers[i].LastScrapeTime = time.Now()
			b.tickers[i].printTicker()
			b.tickers[i].computeHourlySentiment()
			b.tickers[i].pushToDb(d)
			fmt.Printf("Sentiment: %f\n", float64(b.tickers[i].HourlySentiment))
			b.tickers[i].dump_raw()
			b.tickers[i].dump_text()
			b.tickers[i].hourlyWipe()
		}
		time.Sleep(b.interval)
	}
}
