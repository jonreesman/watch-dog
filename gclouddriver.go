package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	d.URI = fmt.Sprintf("%s:%s@unix(/Users/jon/cloudsql/%s)/%s", d.dbUser, d.dbPwd, d.instanceConnection, d.dbName)
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
	_, err := d.db.Exec("DROP TABLE IF EXISTS " + s)
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createTickerTable() {
	_, err := d.db.Exec("CREATE TABLE tickers(ticker_id SERIAL, name VARCHAR(255), num_tweets INTEGER, hourly_sentiment DOUBLE)")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) createStatementTable() {
	_, err := d.db.Exec("CREATE TABLE statements(statement_id SERIAL, expression VARCHAR(255), time_stamp BIGINT, polarity TINYINT)")
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) addTicker(name string, numTweets int, hourlySentiment float64) {
	_, err := d.db.Query(fmt.Sprintf("INSERT INTO tickers(name, num_tweets, hourly_sentiment) VALUES ('%s', '%d', '%f')",
		name,
		numTweets,
		hourlySentiment,
	))
	if err != nil {
		log.Fatal(err)
	}
}

func (d DBManager) addStatement(expression string, timeStamp int64, polarity uint8) {

	_, err := d.db.Query(fmt.Sprintf("INSERT INTO statements(expression, time_stamp, polarity) VALUES (\"%s\", '%d', '%d')",
		expression,
		timeStamp,
		polarity,
	))
	if err != nil {
		log.Fatal(err)
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
