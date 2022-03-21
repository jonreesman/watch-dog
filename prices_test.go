package main

import (
	"fmt"
	"testing"
)

func TestPriceCheck(t *testing.T) {
	fmt.Println("AAPL ", priceCheck("AAPL"))
	fmt.Println("BTC-USD ", priceCheck("BTC-USD"))
}

func TestCheckTickerExists(t *testing.T) {
	fmt.Println("AAPL ", CheckTickerExists("AAPL"))
	fmt.Println("BYXUTK ", CheckTickerExists("BYXUTK"))
}
