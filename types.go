package main

import "time"

type bot struct {
	tickers []ticker
	interval time.Duration //defined in seconds
}

type ticker struct {
	name            string
	lastScrapeTime	time.Time
	numTweets       int
	tweets          []statement
	hourlySentiment float64
}
type statement struct {
	expression string 	`dynamodbav="expression"`
	subject string		`dynamodbav="subject"`
	source string		`dynamodbav="source"`
	timeStamp int64		`dynamodbav="timestamp"`
	polarity   uint8	`dynamodbav="polarity"`
}
