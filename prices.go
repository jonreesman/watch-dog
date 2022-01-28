package main

import (
	"fmt"
	"log"
	"time"

	"github.com/piquette/finance-go/quote"
)

//Intervals of 5 minutes?
func FiveMinutePriceCheck(ticker string) intervalQuote {
	q, err := quote.Get(ticker)
	if err != nil {
		log.Panic(err)
	}
	var iq intervalQuote
	iq.timeObj = time.Now()
	iq.TimeStamp = iq.timeObj.Unix()
	iq.TimeString = iq.timeObj.String()

	iq.CurrentPrice = q.RegularMarketPrice

	return iq
}

func CheckTickerExists(ticker string) bool {
	q, err := quote.Get(ticker)
	fmt.Println(q)
	if err != nil || q == nil {
		return false
	} else {
		return true
	}
}
