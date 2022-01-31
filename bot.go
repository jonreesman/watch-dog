package main

import (
	"strconv"
	"time"

	sentiment "github.com/cdipaolo/sentiment"
)

func initBot() bot {
	var b bot
	b.mainInterval = 3600 * time.Second
	b.quoteInterval = 300 * time.Second
	return b
}

func (b *bot) grabQuotes(d DBManager) { //const d *DBManager
	for {
		for i := range b.tickers {
			j := FiveMinutePriceCheck(b.tickers[i].Name)
			d.addQuote(j.TimeStamp, b.tickers[i].Id, j.CurrentPrice)
		}
		time.Sleep(b.quoteInterval)
	}
}

func (b *bot) addTicker(d DBManager, addTicker chan string) {
	for {
		name := <-addTicker
		t, err := b.tickers.addTicker(name, d)
		if err != nil {
			addTicker <- err.Error()
		} else {
			addTicker <- strconv.Itoa(t.Id)
		}
	}
}

func (b *bot) deleteTicker(d DBManager, deleteTicker chan int) {
	for {
		id := <-deleteTicker
		b.tickers.deleteTicker(id, d)
		deleteTicker <- 200
		/*err := d.deleteTicker(id)
		if err != nil {
			deleteTicker <- 400
		} else {
			deleteTicker <- 200
		}*/
	}
}

func (b bot) run() {
	//MySQL DB Set Up
	var d DBManager
	d.initializeManager()

	//Import tickers from database
	b.tickers.importTickers(d)

	//Prepare channels for API to interface with
	addTicker := make(chan string)
	deleteTicker := make(chan int)

	//Spin off goroutines to handle API CRUD operations
	go b.addTicker(d, addTicker)
	go b.deleteTicker(d, deleteTicker)

	//Spin off server instance with add/deleteTicker channels
	var s Server
	go s.startServer(d, addTicker, deleteTicker)

	//Go routine for collecting market prices every five minutes.
	go b.grabQuotes(d)

	//Instantiate sentiment model
	sentimentModel, err := sentiment.Restore()

	//If it err'd, panic. No way to recover.
	if err != nil {
		panic(err)
	}
	//Main business logic loop of Bot object.
	for {
		//Scrapes all tickers concurrently.
		b.tickers.scrape(sentimentModel)
		//Once scraped, push all to database.
		go b.tickers.pushToDb(d)
		time.Sleep(b.mainInterval)
	}
}
