package main

import (
	"fmt"
	"testing"
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

	tickers.scrape()
	fmt.Println(tickers)
}
