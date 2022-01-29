package main

import "time"

type bot struct {
	tickers       []ticker
	mainInterval  time.Duration //defined in seconds
	quoteInterval time.Duration
	DEBUG         bool
}

//Defines an object packaged for pushing to a database.
type ticker struct {
	Name            string `json:"name"`
	LastScrapeTime  time.Time
	NumTweets       int
	Tweets          []statement `json:"tweets"`
	HourlySentiment float64
	Quotes          []intervalQuote
	id              int
}
type statement struct {
	Expression   string `json:"expression"`
	Subject      string `json:"subject"`
	Source       string `json:"source"`
	TimeStamp    int64  `json:"timeStamp"`
	TimeString   string `json:"timeString"`
	Polarity     uint8  `json:"polarity"`
	timeStampObj time.Time
	URLs         []string
	PermanentURL string
}

/*	Conveniently, we can use the same object for
 *	both sending quotes to the front end as
 *	well as sending sentiments since they are identical
 *	in variable types.
 */
type intervalQuote struct {
	TimeStamp    int64
	TimeString   string
	timeObj      time.Time
	CurrentPrice float64
}
