package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type tickerSlice []ticker

type bot struct {
	tickers      tickerSlice
	mainInterval time.Duration //defined in seconds
}

type Server struct {
	d            DBManager
	router       *gin.Engine
	addTicker    chan string
	deleteTicker chan int
}

type DBManager struct {
	db     *sql.DB
	dbName string
	dbUser string
	dbPwd  string
	URI    string
}

//Defines an object packaged for pushing to a database.
type ticker struct {
	Name            string `json:"name"`
	LastScrapeTime  time.Time
	numTweets       int
	Tweets          []statement `json:"tweets"`
	HourlySentiment float64
	Id              int
	active          int
}
type statement struct {
	Expression   string
	subject      string
	source       string
	TimeStamp    int64
	timeString   string
	Polarity     float64
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
	CurrentPrice float64
}
