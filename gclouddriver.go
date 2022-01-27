package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
	//defer d.db.Close()

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
}

func (d DBManager) dropTable(s string) {
	//REMOVE ONCE DONE DEBUGGING
	_, err := d.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	_, err = d.db.Exec("DROP TABLE IF EXISTS " + s)
	_, err = d.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createTickerTable() {
	_, err := d.db.Exec("CREATE TABLE tickers(ticker_id SERIAL PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createStatementTable() {
	_, err := d.db.Exec("CREATE TABLE statements(statement_id SERIAL PRIMARY KEY, ticker_id BIGINT UNSIGNED, expression VARCHAR(500), url VARCHAR(255), time_stamp BIGINT, polarity TINYINT, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id))")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createSentimentTable() {
	_, err := d.db.Exec("CREATE TABLE sentiments(sentiment_id SERIAL PRIMARY KEY, time_stamp BIGINT, ticker_id BIGINT UNSIGNED, hourly_sentiment DOUBLE, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id))")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createQuotesTable() {
	_, err := d.db.Exec("CREATE TABLE quotes(quote_id SERIAL PRIMARY KEY, time_stamp BIGINT, ticker_id BIGINT UNSIGNED, price DOUBLE, FOREIGN KEY (ticker_id) REFERENCES tickers(ticker_id))")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) addQuote(time_stamp int64, id int, price float64) {
	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO quotes(time_stamp, ticker_id, price) VALUES ('%d', '%d', '%f')",
		time_stamp,
		id,
		price,
	))
	if err != nil {
		log.Print("Error in AddQuote ", err)
	}
}

func (d DBManager) addSentiment(time_stamp int64, id int, hourly_sentiment float64) {
	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO sentiments(time_stamp, ticker_id, hourly_sentiment) VALUES ('%d', '%d', '%f')",
		time_stamp,
		id,
		hourly_sentiment,
	))
	if err != nil {
		log.Print("Error in addSentiment()", err)
	}
}

func (d DBManager) addTicker(name string) int {
	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO tickers(name) VALUES ('%s')",
		name,
	))
	if err != nil {
		log.Print("Error in AddTicker()", err)
	}
	return d.retrieveTicker(name)
}

func (d DBManager) addStatement(expression string, timeStamp int64, polarity uint8, url string) {

	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO statements(expression, time_stamp, polarity, url) VALUES (\"%s\", '%d', '%d', \"%s\")",
		expression,
		timeStamp,
		polarity,
		url,
	))
	if err != nil {
		log.Print("Error in addStatement", err)
	}
}

/*func (d dbManager) addStatements(t ticker) {
	dbQuery, err := d.db.Prepare("INSERT INTO statements(expression, time_stamp, polarity) VALUES (...)")
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range t.Tweets {
		_, err := dbQuery.Query(s.Expression, s.TimeStamp, s.Polarity)
		if err != nil {
			log.Fatal(err)
		}
	}
}*/

func (d DBManager) retrieveTickers() {
	rows, err := d.db.Query("SELECT ticker_id, name FROM tickers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id   string
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v: %s\n", id, name)
	}
}

func (d DBManager) retrieveTicker(tickerName string) int {
	rows, err := d.db.Query("SELECT ticker_id, name FROM tickers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id   string
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		if name == tickerName {
			id_number, _ := strconv.Atoi(id)
			return id_number
		}
	}
	return 0
}

func (d DBManager) returnQuoteHistory(id int) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT time_stamp, price FROM quotes WHERE ticker_id=%d ORDER BY time_stamp", id))
	if err != nil {
		log.Print("Error returning Quote History: ", err)
	}
	//fmt.Print(rows)
	var iq []intervalQuote
	var q intervalQuote
	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&q.TimeStamp, &q.currentPrice)
		iq = append(iq, q)
		if err != nil {
			log.Print("Error in row scan")
		}
	}
}

func (d DBManager) returnSentimentHistory(id int) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT time_stamp, hourly_sentiment FROM sentiments WHERE ticker_id=%d ORDER BY time_stamp", id))
	if err != nil {
		log.Print("Error returning senitment history: ", err)
	}

	var sh []intervalQuote
	var s intervalQuote

	for rows.Next() {
		if rows.Err() != nil {
			log.Print("Found no rows.")
		}
		err := rows.Scan(&s.TimeStamp, &s.currentPrice)
		sh = append(sh, s)
		if err != nil {
			log.Print("Error in row scan")
		}
	}

	for i, s := range sh {
		fmt.Println(i, s)
	}
}
