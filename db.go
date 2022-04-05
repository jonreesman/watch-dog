package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBManager struct {
	db     *sql.DB
	dbName string
	dbUser string
	dbPwd  string
	URI    string
}

func (d *DBManager) initializeManager() error {
	d.dbUser = os.Getenv("DB_USER")
	d.dbPwd = os.Getenv("DB_PWD")
	d.dbName = os.Getenv("DB_NAME")

	d.URI = fmt.Sprintf("%s:%s@tcp(%s)/%s", d.dbUser, d.dbPwd, "127.0.0.1:3306", d.dbName)
	d.db, _ = sql.Open("mysql", d.URI)

	err := d.db.Ping()
	if err != nil {
		return err
	}
	_, err = d.db.Exec(fmt.Sprintf("USE %s", d.dbName))
	if err != nil {
		return err
	}
	fmt.Println("Connection established")

	//d.dropTable("tickers")
	//d.dropTable("statements")
	//d.dropTable("sentiments")
	d.createTickerTable()
	d.createStatementTable()
	d.createSentimentTable()
	return nil
}

const addSentimentQuery = `
INSERT INTO sentiments(time_stamp, ticker_id, hourly_sentiment) ` +
	`VALUES (?, ?, ?)`

func (d DBManager) addSentiment(wg *sync.WaitGroup, timeStamp int64, tickerId int, hourlySentiment float64) {
	defer wg.Done()
	_, err := d.db.Exec(addSentimentQuery,
		timeStamp,
		tickerId,
		float32(hourlySentiment),
	)
	if err != nil {
		log.Print("Error in addSentiment()", err, hourlySentiment)
		log.Print("id is ", tickerId)
	}
}

const updateTickerQuery = `
UPDATE tickers SET last_scrape_time=? ` +
	`WHERE ticker_id=?`

func (d DBManager) updateTicker(wg *sync.WaitGroup, id int, timeStamp time.Time) error {
	defer wg.Done()
	if _, err := d.db.Exec(updateTickerQuery, timeStamp.Unix(), id); err != nil {
		return err
	}
	return nil
}

const activateTickerQuery = `
UPDATE tickers SET active=1 WHERE ticker_id=?`

const addTickerQuery = `
INSERT INTO tickers(name, active, last_scrape_time) ` +
	`VALUES (?,?,?)`

func (d DBManager) addTicker(name string) (int, error) {
	if t, err := d.retrieveTickerByName(name); err == nil {
		if t.Active == 1 {
			return t.Id, errors.New("ticker already exists and is active")
		}
		if _, err := d.db.Exec(activateTickerQuery, t.Id); err != nil {
			return 0, err
		}
		return t.Id, nil
	}

	dbQuery, err := d.db.Prepare(addTickerQuery)
	if err != nil {
		log.Print("Error in AddTicker()", err)
		return 0, errors.New("failed to add ticker")
	}
	_, err = dbQuery.Query(name, 1, nil)
	if err != nil {
		log.Print("Error in AddTicker()", err)
		return 0, errors.New("failed to add ticker")
	}
	var t ticker
	if t, err = d.retrieveTickerByName(name); err != nil {
		log.Printf("Failed to add ticker")
		return 0, errors.New("failed to add ticker")
	}

	return t.Id, nil

}

const addStatementQuery = `
INSERT INTO statements(ticker_id, expression, time_stamp, polarity, url, tweet_id) ` +
	`VALUES (?, ?, ?, ?, ?, ?)`

func (d DBManager) addStatement(wg *sync.WaitGroup, tickerId int, expression string, timeStamp int64, polarity float64, url string, tweet_id uint64) {
	defer wg.Done()
	_, err := d.db.Exec(addStatementQuery,
		tickerId,
		expression,
		timeStamp,
		float32(polarity),
		url,
		tweet_id,
	)
	if err != nil {
		log.Print("Error in addStatement", err)
	}
}

const activeTickerQuery = `
SELECT tickers.ticker_id, tickers.name, tickers.last_scrape_time, ` +
	`sentiments.hourly_sentiment FROM tickers LEFT JOIN sentiments ` +
	`ON tickers.ticker_id = sentiments.ticker_id ` +
	`AND tickers.last_scrape_time = sentiments.time_stamp ` +
	`WHERE active=1 ORDER BY ticker_id`

func (d DBManager) returnActiveTickers() (tickers tickerSlice) {
	rows, err := d.db.Query(activeTickerQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id                    int
		name                  string
		lastScrapeTimeHolder  sql.NullInt64
		lastScrapeTime        int64
		hourlySentimentHolder sql.NullFloat64
		hourlySentiment       float64
	)

	for rows.Next() {
		err := rows.Scan(&id, &name, &lastScrapeTimeHolder, &hourlySentimentHolder)
		if err != nil {
			log.Print(err)
		}
		if lastScrapeTimeHolder.Valid {
			lastScrapeTime = lastScrapeTimeHolder.Int64
		} else {
			lastScrapeTime = 0
		}
		if hourlySentimentHolder.Valid {
			hourlySentiment = hourlySentimentHolder.Float64
		} else {
			hourlySentiment = 0
		}
		tickers.appendTicker(ticker{
			Name:            name,
			Id:              id,
			LastScrapeTime:  time.Unix(lastScrapeTime, 0),
			HourlySentiment: hourlySentiment,
		})
		log.Printf("%v: %s\n", id, name)
	}
	return tickers
}

const returnAllTickersQuery = `
SELECT ticker_id, name, last_scrape_time ` +
	`FROM tickers ORDER BY ticker_id`

func (d DBManager) returnAllTickers() (tickers tickerSlice) {
	rows, err := d.db.Query(returnAllTickersQuery)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	var (
		id             int
		name           string
		lastScrapeTime sql.NullInt64
	)

	for rows.Next() {
		err := rows.Scan(&id, &name, &lastScrapeTime)
		if err != nil {
			log.Fatal(err)
		}
		tickers.appendTicker(ticker{Name: name, Id: id, LastScrapeTime: time.Unix(lastScrapeTime.Int64, 0)})
		log.Printf("%v: %s\n", id, name)
	}
	return tickers
}

const retrieveTickerByNameQuery = `
SELECT ticker_id, name, last_scrape_time, active FROM tickers ` +
	`WHERE name=?`

func (d DBManager) retrieveTickerByName(tickerName string) (ticker, error) {
	rows, err := d.db.Query(retrieveTickerByNameQuery, tickerName)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	var (
		id             int
		name           string
		lastScrapeTime sql.NullInt64
		active         int
	)
	for rows.Next() {
		err := rows.Scan(&id, &name, &lastScrapeTime, &active)
		if err != nil {
			log.Print(err)
		}
		if name == tickerName {
			t := ticker{
				Id:             id,
				Name:           name,
				LastScrapeTime: time.Unix(lastScrapeTime.Int64, 0),
				Active:         active,
			}
			return t, nil
		}
	}
	return ticker{Id: 0, Name: ""}, errors.New("ticker does not exist with that ID")
}

const retrieveTickerByIdQuery = `
SELECT ticker_id, name, last_scrape_time FROM tickers ` +
	`WHERE ticker_id=?`

func (d DBManager) retrieveTickerById(tickerId int) (ticker, error) {
	rows, err := d.db.Query(retrieveTickerByIdQuery, tickerId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id             string
		name           string
		lastScrapeTime sql.NullInt64
	)
	strId := strconv.Itoa(tickerId)
	for rows.Next() {
		err := rows.Scan(&id, &name, &lastScrapeTime)
		if err != nil {
			log.Fatal(err)
		}
		if strId == id {
			return ticker{Name: name, Id: tickerId, LastScrapeTime: time.Unix(lastScrapeTime.Int64, 0)}, nil
		}
	}
	return ticker{Name: "none"}, errors.New("ticker does not exist")
}

const returnSentimentHistoryQuery = `
SELECT time_stamp, hourly_sentiment FROM sentiments ` +
	`WHERE ticker_id=? ORDER BY time_stamp DESC`

func (d DBManager) returnSentimentHistory(id int, fromTime int64) []intervalQuote {
	rows, err := d.db.Query(returnSentimentHistoryQuery, id)
	if err != nil {
		log.Print("Error returning senitment history: ", err)
	}

	var payload []intervalQuote
	var s intervalQuote

	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&s.TimeStamp, &s.CurrentPrice)
		if s.TimeStamp < fromTime {
			break
		}
		payload = append(payload, s)
		if err != nil {
			log.Print("Error in row scan")
		}
	}
	return payload
}

const returnAllStatementsQuery = `
SELECT time_stamp, expression, url, polarity, tweet_id ` +
	`FROM statements WHERE ticker_id=? ` +
	`ORDER BY time_stamp DESC`

func (d DBManager) returnAllStatements(id int, fromTime int64) []statement {
	rows, err := d.db.Query(returnAllStatementsQuery, id)
	if err != nil {
		log.Print("Error returning senitment history: ", err)
	}

	var (
		returnPackage []statement
		st            statement
	)

	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&st.TimeStamp, &st.Expression, &st.PermanentURL, &st.Polarity, &st.ID)
		if st.TimeStamp < fromTime {
			break
		}
		returnPackage = append(returnPackage, st)
		if err != nil {
			log.Print("Error in row scan")
		}
	}
	return returnPackage
}

const deactivateTickerQuery = `
UPDATE tickers SET active=0 WHERE ticker_id=?`

func (d DBManager) deactivateTicker(id int) error {
	if _, err := d.db.Exec(deactivateTickerQuery, id); err != nil {
		return err
	}
	return nil
}
