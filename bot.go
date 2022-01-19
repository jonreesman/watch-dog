package main

import (
	"fmt"
	"time"
)

func initBot() bot {
	var b bot
	b.interval = 3600 * time.Second
	b.tickers = importTickers()
	return b
}

func (b bot) run() {

	//MySQL DB Set Up
	/*var d DBManager
	d.initializeManager()
	d.dropTable("tickers")
	d.dropTable("statements")
	d.createTickerTable()
	d.createStatementTable()
	for _, stock := range b.tickers {
		d.addTicker(stock.Name, stock.NumTweets, stock.HourlySentiment)
		fmt.Println(stock)
	}
	d.retrieveTickers()*/

	//Main business logic loop of Bot object.
	for {
		fmt.Println("Loop")
		scrapeAll(&b.tickers)
		for i := range b.tickers {
			b.tickers[i].LastScrapeTime = time.Now()
			b.tickers[i].printTicker()
			b.tickers[i].computeHourlySentiment()
			//b.tickers[i].pushToDb(d)

			//
			b.tickers[i].dump_raw()
			b.tickers[i].dump_text()

			//We do not keep the tweets cached hour to hour,
			//so we wipe them since they are accesible in the database.
			b.tickers[i].hourlyWipe()
		}
		time.Sleep(b.interval)
	}
}
