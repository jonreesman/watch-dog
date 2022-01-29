package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sentiment "github.com/cdipaolo/sentiment"
)

func initBot() bot {
	var b bot
	b.DEBUG, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	b.mainInterval = 3600 * time.Second
	b.quoteInterval = 300 * time.Second
	return b
}

func (b *bot) grabQuotes(d DBManager) { //const d *DBManager
	for {
		for i := range b.tickers {
			j := FiveMinutePriceCheck(b.tickers[i].Name)
			b.tickers[i].Quotes = append(b.tickers[i].Quotes, j)
			if b.DEBUG {
				fmt.Println(b.tickers[i].Name, ":", j, "id:", b.tickers[i].id)
			}
			d.addQuote(j.TimeStamp, b.tickers[i].id, j.CurrentPrice)
		}
		time.Sleep(b.quoteInterval)
	}
}

func (b bot) run() {
	//MySQL DB Set Up
	var d DBManager
	d.initializeManager()
	b.tickers = importTickers(d)

	var s Server
	go s.startServer(d, &b.tickers)

	go b.grabQuotes(d)
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}
	//Main business logic loop of Bot object.
	for {
		if b.DEBUG {
			fmt.Println("Loop")
		}
		scrapeAll(&b.tickers)
		for i := range b.tickers {
			if i == 0 {
				continue
			}
			b.tickers[i].LastScrapeTime = time.Now()
			b.tickers[i].computeHourlySentiment(sentimentModel)
			b.tickers[i].pushToDb(d)

			//We do not keep the tweets cached hour to hour,
			//so we wipe them since they are accesible in the database.
			b.tickers[i].hourlyWipe()
		}
		time.Sleep(b.mainInterval)
	}
}
