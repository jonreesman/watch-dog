package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTwitterScrape(t *testing.T) {
	tick := ticker{
		Name:            "AMD",
		numTweets:       0,
		Tweets:          []statement{},
		HourlySentiment: 0,
		Id:              0,
		Active:          0,
	}
	statements := twitterScrape(tick)
	var maxTime int64
	minTime := time.Now().Unix()
	for _, s := range statements {
		if s.TimeStamp < minTime {
			minTime = s.TimeStamp
		}
		if s.TimeStamp > maxTime {
			maxTime = s.TimeStamp
		}
	}
	fmt.Println("Number of tweets: " + strconv.Itoa(len(statements)))
	fmt.Println("minTime: " + time.Unix(minTime, 0).String() + " maxTime: " + time.Unix(maxTime, 0).String())
}
