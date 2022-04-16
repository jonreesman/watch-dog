package main

import (
	"fmt"
	"testing"
)

func TestPriceCheck(t *testing.T) {
	fmt.Println("AAPL ", priceCheck("AAPL"))
	fmt.Println("BTC ", priceCheck("BTC"))

}

func TestCheckTickerExists(t *testing.T) {
	fmt.Println("AAPL ", CheckTickerExists("AAPL"))
	fmt.Println("BYXUTK ", CheckTickerExists("BYXUTK"))
}
