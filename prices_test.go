package main

import (
	"fmt"
	"testing"
)

func TestFiveMinutePriceCheck(t *testing.T) {
	fmt.Println("AAPL ", FiveMinutePriceCheck("AAPL"))
	fmt.Println("BTC-USD ", FiveMinutePriceCheck("BTC-USD"))
}

func TestCheckTickerExists(t *testing.T) {
	fmt.Println("AAPL ", CheckTickerExists("AAPL"))
	fmt.Println("BYXUTK ", CheckTickerExists("BYXUTK"))
}
