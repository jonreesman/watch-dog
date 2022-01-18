package main

import "time"

type bot struct {
	tickers  []ticker
	interval time.Duration //defined in seconds
}

//Defines an object packaged for pushing to a database.
type ticker struct {
	Name            string      `json:"name"`
	LastScrapeTime  time.Time   `json:"lastScrapeTime"`
	NumTweets       int         `json:"numTweets"`
	Tweets          []statement `json:"tweets"`
	HourlySentiment float64     `json:"hourlySentiment"`
}
type statement struct {
	Expression   string `json:"expression" dynamodbav:"expression" bson:"expression"`
	Subject      string `json:"subject" dynamodbav:"subject" bson:"subject"`
	Source       string `json:"source" dynamodbav:"source" bson:"source"`
	TimeStamp    int64  `json:"timeStamp" dynamodbav:"timestamp" bson:"timestamp"`
	TimeString   string `json:"timeString" dynamodbav:"timeString" bson:"timeString"`
	Polarity     uint8  `json:"polarity" dynamodbav:"polarity" bson:"polarity"`
	timeStampObj time.Time
}
