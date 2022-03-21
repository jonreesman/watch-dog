package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/jonreesman/watch-dog/pb"
)

func (b *bot) initBot(d DBManager) error {
	b.mainInterval = 3600 * time.Second
	if err := b.tickers.importTickers(d); err != nil {
		log.Printf("initBot(): Failed to import tickers.")
		return err
	}
	return nil
}

func AddTicker(d DBManager, addTicker chan string) {
	for {
		name := <-addTicker
		s := sanitize(name)
		if !CheckTickerExists(s) {
			log.Println("Stock does not exist.")

			addTicker <- errors.New("stock/crypto does not exist").Error()
		}
		t, err := d.retrieveTickerByName(name)
		if err == nil {
			if t.active == 1 {
				addTicker <- strconv.Itoa(t.Id)
				continue
			} else {
				t.active = 1
			}
		}
		if err != nil {
			t = ticker{
				Name:            s,
				LastScrapeTime:  time.Time{},
				numTweets:       0,
				Tweets:          []statement{},
				HourlySentiment: 0,
				Id:              t.Id,
				active:          1,
			}
		}

		if id, err := d.addTicker(s); err != nil {
			addTicker <- err.Error()
		} else {
			t.Id = id
		}
		var wg sync.WaitGroup
		wg.Add(1)
		t.scrape(&wg)
		//j := FiveMinutePriceCheck(t.Name)
		//d.addQuote(j.TimeStamp, t.Id, j.CurrentPrice)
		wg.Wait()
		t.pushToDb(d)

		addTicker <- strconv.Itoa(t.Id)
	}
}

func DeactivateTicker(d DBManager, deleteTicker chan int) {
	for {
		id := <-deleteTicker
		err := d.deactivateTicker(id)
		if err != nil {
			deleteTicker <- 400
		}
		deleteTicker <- 200
	}
}

func (b bot) run(d DBManager) {

	//Main business logic loop of Bot object.
	for {
		//Scrapes all tickers concurrently.
		b.tickers.scrape()
		//Once scraped, push all to database.
		go b.tickers.pushToDb(d)
		time.Sleep(b.mainInterval)
		b.tickers = nil
		b.tickers.importTickers(d)
	}
}
