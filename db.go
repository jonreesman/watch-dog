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

	d.dropTable("tickers")
	d.dropTable("statements")
	//d.dropTable("quotes")
	d.dropTable("sentiments")
	d.createTickerTable()
	d.createStatementTable()
	//d.createQuotesTable()
	d.createSentimentTable()
	return nil
}

func (d DBManager) addQuote(timeStamp int64, id int, price float64) {
	_, err := d.db.Exec("INSERT INTO quotes(time_stamp, ticker_id, price) VALUES (?, ?, ?)",
		timeStamp,
		id,
		price,
	)
	if err != nil {
		log.Print("Error in AddQuote ", err)
		log.Print("id is ", id)
	}
}

func (d DBManager) addSentiment(wg *sync.WaitGroup, timeStamp int64, tickerId int, hourlySentiment float64) {
	defer wg.Done()
	_, err := d.db.Exec("INSERT INTO sentiments(time_stamp, ticker_id, hourly_sentiment) VALUES (?, ?, ?)",
		timeStamp,
		tickerId,
		float32(hourlySentiment),
	)
	if err != nil {
		log.Print("Error in addSentiment()", err, hourlySentiment)
		log.Print("id is ", tickerId)
	}
}

func (d DBManager) updateTicker(wg *sync.WaitGroup, id int, timeStamp time.Time) error {
	defer wg.Done()
	if _, err := d.db.Exec("UPDATE tickers SET last_scrape_time=? WHERE ticker_id=?", timeStamp.Unix(), id); err != nil {
		return err
	}
	return nil
}

func (d DBManager) addTicker(name string) (int, error) {
	if t, err := d.retrieveTickerByName(name); err == nil {
		if t.active == 1 {
			return t.Id, errors.New("ticker already exists and is active")
		}
		if _, err := d.db.Exec("UPDATE tickers SET active=1 WHERE ticker_id=?", t.Id); err != nil {
			return 0, err
		}
		return t.Id, nil
	}

	dbQuery, err := d.db.Prepare("INSERT INTO tickers(name, active, last_scrape_time) VALUES (?,?,?)")
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

func (d DBManager) addStatement(wg *sync.WaitGroup, tickerId int, expression string, timeStamp int64, polarity float64, url string) {
	defer wg.Done()
	_, err := d.db.Exec("INSERT INTO statements(ticker_id, expression, time_stamp, polarity, url) VALUES (?, ?, ?, ?, ?)",
		tickerId,
		expression,
		timeStamp,
		float32(polarity),
		url,
	)
	if err != nil {
		log.Print("Error in addStatement", err)
	}
}

func (d DBManager) returnActiveTickers() (tickers tickerSlice) {
	//Change to do a union to return last sentiment

	rows, err := d.db.Query("SELECT tickers.ticker_id, tickers.name, tickers.last_scrape_time, sentiments.hourly_sentiment FROM tickers LEFT JOIN sentiments ON tickers.ticker_id = sentiments.ticker_id WHERE active=1 ORDER BY ticker_id")
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

func (d DBManager) returnAllTickers() (tickers tickerSlice) {
	rows, err := d.db.Query("SELECT ticker_id, name, last_scrape_time FROM tickers ORDER BY ticker_id")
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

func (d DBManager) retrieveTickerByName(tickerName string) (ticker, error) {
	rows, err := d.db.Query("SELECT ticker_id, name, last_scrape_time, active FROM tickers")
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
				active:         active,
			}
			return t, nil
		}
	}
	return ticker{Id: 0, Name: ""}, errors.New("ticker does not exist with that ID")
}

func (d DBManager) retrieveTickerById(tickerId int) (ticker, error) {
	rows, err := d.db.Query("SELECT ticker_id, name, last_scrape_time FROM tickers")
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

func (d DBManager) returnSentimentHistory(id int, fromTime int64) []intervalQuote {
	rows, err := d.db.Query("SELECT time_stamp, hourly_sentiment FROM sentiments WHERE ticker_id=? ORDER BY time_stamp DESC", id)
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

func (d DBManager) returnAllStatements(id int, fromTime int64) []statement {
	rows, err := d.db.Query("SELECT time_stamp, expression, url, polarity FROM statements WHERE ticker_id=? ORDER BY time_stamp DESC", id)
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
		err := rows.Scan(&st.TimeStamp, &st.Expression, &st.PermanentURL, &st.Polarity)
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

func (d DBManager) deactivateTicker(id int) error {
	if _, err := d.db.Exec("UPDATE tickers SET active=0 WHERE ticker_id=?", id); err != nil {
		return err
	}
	return nil
}
