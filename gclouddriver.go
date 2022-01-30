package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DBManager struct {
	db                 *sql.DB
	dbName             string
	dbUser             string
	dbPwd              string
	instanceConnection string
	URI                string
}

func (d *DBManager) initializeManager() {
	d.dbUser = os.Getenv("DB_USER")
	d.dbPwd = os.Getenv("DB_PWD")
	d.instanceConnection = os.Getenv("INSTANCE_CONNECTION_NAME")
	d.dbName = os.Getenv("DB_NAME")

	//d.URI = fmt.Sprintf("%s:%s@unix(/Users/jon/cloudsql/%s)/%s", d.dbUser, d.dbPwd, d.instanceConnection, d.dbName)
	d.URI = fmt.Sprintf("%s:%s@tcp(%s)/%s", d.dbUser, d.dbPwd, "127.0.0.1:1433", d.dbName)
	d.db, _ = sql.Open("mysql", d.URI)

	err := d.db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("We got here...")
	_, err = d.db.Exec(fmt.Sprintf("USE %s", d.dbName))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established")

	//d.dropTable("tickers")
	//d.dropTable("statements")
	//d.dropTable("quotes")
	//d.dropTable("sentiments")
	d.createTickerTable()
	d.createStatementTable()
	d.createQuotesTable()
	d.createSentimentTable()
}

func (d DBManager) dropTable(s string) {
	//REMOVE ONCE DONE DEBUGGING
	_, err := d.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.db.Exec("DROP TABLE IF EXISTS " + s)
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createTickerTable() {
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS tickers(ticker_id SERIAL PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.db.Exec("ALTER TABLE tickers ADD CONSTRAINT ticker_Unique UNIQUE(name)")
	if err != nil {
		log.Print(err)
	}
}

func (d DBManager) createStatementTable() {
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS statements(statement_id SERIAL PRIMARY KEY, ticker_id BIGINT UNSIGNED, expression VARCHAR(500), url VARCHAR(255), time_stamp BIGINT, polarity TINYINT, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id) ON DELETE CASCADE)")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createSentimentTable() {
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS sentiments(sentiment_id SERIAL PRIMARY KEY, time_stamp BIGINT, ticker_id BIGINT UNSIGNED, hourly_sentiment FLOAT, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id) ON DELETE CASCADE)")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createQuotesTable() {
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS quotes(quote_id SERIAL PRIMARY KEY, time_stamp BIGINT, ticker_id BIGINT UNSIGNED, price DOUBLE, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id) ON DELETE CASCADE)")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) addQuote(time_stamp int64, id int, price float64) {
	_, err := d.db.Exec("INSERT INTO quotes(time_stamp, ticker_id, price) VALUES (?, ?, ?)",
		time_stamp,
		id,
		price,
	)
	if err != nil {
		log.Print("Error in AddQuote ", err)
		log.Print("id is ", id)
	}
}

func (d DBManager) addSentiment(wg *sync.WaitGroup, time_stamp int64, id int, hourly_sentiment float64) {
	defer wg.Done()
	_, err := d.db.Exec("INSERT INTO sentiments(time_stamp, ticker_id, hourly_sentiment) VALUES (?, ?, ?)",
		time_stamp,
		id,
		float32(hourly_sentiment),
	)
	if err != nil {
		log.Print("Error in addSentiment()", err, hourly_sentiment)
		log.Print("id is ", id)
	}
}

func (d DBManager) addTicker(name string) (int, error) {
	dbQuery, err := d.db.Prepare("INSERT INTO tickers(name) VALUES (?)")
	if err != nil {
		log.Print("Error in AddTicker()", err)
		return 0, errors.New("failed to add ticker")
	}
	_, err = dbQuery.Query(name)
	if err != nil {
		log.Print("Error in AddTicker()", err)
		return 0, errors.New("failed to add ticker")
	}
	var id int
	if id, err = d.retrieveTickerByName(name); err != nil {
		log.Printf("Failed to add ticker")
		return 0, errors.New("failed to add ticker")
	}

	return id, nil

}

func (d DBManager) addStatement(wg *sync.WaitGroup, expression string, timeStamp int64, polarity uint8, url string) {
	defer wg.Done()
	_, err := d.db.Exec("INSERT INTO statements(expression, time_stamp, polarity, url) VALUES (?, ?, ?, ?)",
		expression,
		timeStamp,
		polarity,
		url,
	)
	if err != nil {
		log.Print("Error in addStatement", err)
	}
}

func (d DBManager) returnTickers() (tickers tickerSlice) {
	rows, err := d.db.Query("SELECT ticker_id, name FROM tickers ORDER BY ticker_id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id   int
		name string
	)

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		tickers.appendTicker(ticker{Name: name, Id: id})
		log.Printf("%v: %s\n", id, name)
	}
	return tickers
}

func (d DBManager) retrieveTickerByName(tickerName string) (int, error) {
	rows, err := d.db.Query("SELECT ticker_id, name FROM tickers")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	var (
		id   int
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Print(err)
		}
		if name == tickerName {
			return id, nil
		}
	}
	return 0, errors.New("ticker does not exist with that ID")
}

func (d DBManager) retrieveTickerById(tickerId int) (ticker, error) {
	rows, err := d.db.Query("SELECT ticker_id, name FROM tickers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id   string
		name string
	)
	strId := strconv.Itoa(tickerId)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		if strId == id {
			return ticker{Name: name, Id: tickerId}, nil
		}
	}
	return ticker{Name: "none"}, errors.New("ticker does not exist")
}

func (d DBManager) returnQuoteHistory(id int, fromTime int64) []intervalQuote {
	rows, err := d.db.Query("SELECT time_stamp, price FROM quotes WHERE ticker_id=? ORDER BY time_stamp DESC", id)
	if err != nil {
		log.Print("Error returning Quote History: ", err)
	}
	var iq []intervalQuote
	var q intervalQuote
	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&q.TimeStamp, &q.CurrentPrice)
		if q.TimeStamp < fromTime {
			break
		}
		iq = append(iq, q)
		if err != nil {
			log.Print("Error in row scan")
		}
	}
	return iq
}

func (d DBManager) returnSentimentHistory(id int, fromTime int64) []intervalQuote {
	rows, err := d.db.Query("SELECT time_stamp, hourly_sentiment FROM sentiments WHERE ticker_id=? ORDER BY time_stamp DESC", id)
	if err != nil {
		log.Print("Error returning senitment history: ", err)
	}

	var sh []intervalQuote
	var s intervalQuote

	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&s.TimeStamp, &s.CurrentPrice)
		if s.TimeStamp < fromTime {
			break
		}
		sh = append(sh, s)
		if err != nil {
			log.Print("Error in row scan")
		}
	}
	return sh
}

func (d DBManager) deleteTicker(id int) error {
	if _, err := d.db.Exec("DELETE FROM tickers WHERE ticker_id=?", id); err != nil {
		return err
	}
	return nil
}
