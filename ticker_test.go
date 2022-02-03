package main

import (
	"log"
	"testing"

	"github.com/cdipaolo/sentiment"
)

func TestScrape(t *testing.T) {
	tickers := make(tickerSlice, 10)
	tickers[0].Name = "AMD"
	tickers[1].Name = "XLNX"
	tickers[2].Name = "RKLB"
	tickers[3].Name = "TSLA"
	tickers[4].Name = "AAPL"
	tickers[5].Name = "BTC-USD"
	tickers[6].Name = "LCID"
	tickers[7].Name = "AMC"
	tickers[8].Name = "GME"
	tickers[9].Name = "UBER"
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Print(err)
	}
	tickers.scrape(sentimentModel)
	//fmt.Println(tickers)
}
