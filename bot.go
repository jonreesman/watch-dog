package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func initBot() bot {
	var b bot
	b.DEBUG, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	b.mainInterval = 3600 * time.Second
	b.quoteInterval = 300 * time.Second
	b.tickers = importTickers()
	return b
}

func (b *bot) grabQuotes() {
	for {
		for i := range b.tickers {
			b.tickers[i].Quotes = append(b.tickers[i].Quotes, FiveMinutePriceCheck(b.tickers[i].Name))
			if b.DEBUG {
				fmt.Println(b.tickers[i].Name, ":", b.tickers[i].Quotes[len(b.tickers[i].Quotes)-1])
			}
		}
		time.Sleep(b.quoteInterval)
	}
}

func (b bot) run() {
	go b.grabQuotes()
	//MySQL DB Set Up
	/*var d DBManager
	d.initializeManager()
	d.dropTable("tickers")
	d.dropTable("statements")
	d.createTickerTable()
	d.createStatementTable()
	d.createQuotesTable()
	d.createSentimentTable()
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
		time.Sleep(b.mainInterval)
	}
}
