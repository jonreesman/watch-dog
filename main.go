package main

import (
	"log"
	"time"

	"github.com/cdipaolo/sentiment"
)

func main() {
	run()
}

func run() {
	var (
		d DBManager
		b bot
		s Server
	)
	d.initializeManager()
	b.initBot()
	b.tickers.importTickers(d)
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Fatal("Unrecoverable error: ", err)
	}

	addTickerChannel := make(chan string)
	deactivateTickerChannel := make(chan int)

	go s.startServer(d, addTickerChannel, deactivateTickerChannel)
	go AddTicker(d, addTickerChannel, sentimentModel)
	go DeactivateTicker(d, deactivateTickerChannel)
	go grabQuotes(d, 300*time.Second)
	b.run(d, sentimentModel)
}
