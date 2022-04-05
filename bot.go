package main

import (
	"errors"
	"log"
	"sync"
	"time"

	_ "github.com/jonreesman/watch-dog/pb"
)

type bot struct {
	tickers      tickerSlice
	mainInterval time.Duration //defined in seconds
}

func (b *bot) initBot(d DBManager) error {
	b.mainInterval = 3600 * time.Second
	if err := b.tickers.importTickers(d); err != nil {
		log.Printf("initBot(): Failed to import tickers.")
		return err
	}
	return nil
}

func AddTicker(d DBManager, name string) (int, error) {
	s := sanitize(name)
	if !CheckTickerExists(s) {
		log.Println("Stock does not exist.")
		return 0, errors.New("stock/crypto does not exist")
	}
	t, err := d.retrieveTickerByName(name)
	if err != nil {
		t = ticker{
			Name:            s,
			LastScrapeTime:  time.Time{},
			numTweets:       0,
			Tweets:          []statement{},
			HourlySentiment: 0,
			Id:              t.Id,
			Active:          1,
		}
	} else {
		if t.Active == 1 {
			return t.Id, nil
		} else {
			t.Active = 1
		}
	}

	if id, err := d.addTicker(s); err != nil {
		return 0, errors.New("addTicker(): " + err.Error())
	} else {
		t.Id = id
	}
	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		t.scrape(&wg)
		wg.Wait()
		t.pushToDb(d)
	}()
	return t.Id, nil
}

func DeactivateTicker(d DBManager, id int) error {
	err := d.deactivateTicker(id)
	if err != nil {
		return err
	}
	return nil
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
