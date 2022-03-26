package main

import (
	"bufio"
	"context"

	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jonreesman/watch-dog/pb"
	"google.golang.org/grpc"
)

type tickerSlice []ticker

//Defines an object packaged for pushing to a database.
type ticker struct {
	Name            string
	LastScrapeTime  time.Time
	numTweets       int
	Tweets          []statement
	HourlySentiment float64
	Id              int
	Active          int
}

type statement struct {
	Expression   string
	subject      string
	source       string
	TimeStamp    int64
	Polarity     float64
	URLs         []string
	PermanentURL string
	ID           uint64
}

type intervalQuote struct {
	TimeStamp    int64
	CurrentPrice float64
}

func (t *ticker) hourlyWipe() {
	t.numTweets = 0
	t.Tweets = nil

}

func (tickers *tickerSlice) pushToDb(d DBManager) {
	ts := *tickers
	for i := range ts {
		go func(t *ticker) {
			t.pushToDb(d)
			t.hourlyWipe()
		}(&ts[i])
	}
	*tickers = ts
}

func (t ticker) pushToDb(d DBManager) {
	var wg sync.WaitGroup
	for _, tw := range t.Tweets {
		fmt.Println("added statement to DB for:", tw.subject)
		wg.Add(1)
		go d.addStatement(&wg, t.Id, tw.Expression, tw.TimeStamp, tw.Polarity, tw.PermanentURL, tw.ID)
	}
	wg.Add(1)
	go d.addSentiment(&wg, t.LastScrapeTime.Unix(), t.Id, t.HourlySentiment)
	wg.Add(1)
	go d.updateTicker(&wg, t.Id, t.LastScrapeTime)
	wg.Wait()
}

func (tickers *tickerSlice) importTickers(d DBManager) error {
	existingTickers := d.returnActiveTickers()
	ts := *tickers
	if len(existingTickers) != 0 {
		*tickers = existingTickers
		return nil
	}
	file, err := os.Open("tickers.txt")
	if err != nil {
		log.Panicf("importTickers(): failed reading data from file\n")
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tick := ticker{Name: scanner.Text(), numTweets: 0}
		if CheckTickerExists(tick.Name) {
			fmt.Println("Exists")
			tick.Id, _ = d.addTicker(tick.Name)
			ts.appendTicker(tick)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ImportTickers(): Error in scanning file.")
		return err
	}
	*tickers = ts
	return nil
}

func (tickers *tickerSlice) appendTicker(t ticker) {
	ts := *tickers
	ts = append(ts, t)
	*tickers = ts
}

func (tickers *tickerSlice) scrape() {
	var wg sync.WaitGroup
	ts := *tickers
	for i := range ts {
		wg.Add(1)
		go ts[i].scrape(&wg)
	}
	wg.Wait()
	*tickers = ts
}

func (t *ticker) scrape(wg *sync.WaitGroup) {
	defer wg.Done()

	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.numTweets = len(t.Tweets)
	t.LastScrapeTime = time.Now()
	wg.Add(1)
	t.computeHourlySentiment(wg)
}

type pack struct {
	Tweet string `json:"tweet"`
}

func (t *ticker) computeHourlySentiment(wg *sync.WaitGroup) {
	defer wg.Done()
	var total float64
	//var response float64
	addr := "localhost:9999"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("computeHourlySentiment(): Failed to dial GRPC.")
		return
	}
	defer conn.Close()
	client := pb.NewSentimentClient(conn)
	for i, s := range t.Tweets {
		/*p := pack{Tweet: s.Expression}
		js, _ := json.Marshal(p)
		pythonSentiment, err := http.Post("http://localhost:8000/tweet", "application/json", bytes.NewBuffer(js))
		if err != nil {
			log.Print(err)
		}
		resp, err := ioutil.ReadAll(pythonSentiment.Body)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal([]byte(resp), &response)
		t.Tweets[i].Polarity = response*/
		request := pb.SentimentRequest{
			Tweet: s.Expression,
		}

		response, err := client.Detect(context.Background(), &request)
		if err != nil {
			log.Printf("GRPC SentimentRequest: %v", err)
		}
		t.Tweets[i].Polarity = float64(response.Polarity)
		total += float64(response.Polarity)
	}
	t.HourlySentiment = total / float64(t.numTweets)
}

//DEPRECATED
/*func (t *ticker) singleScrape() {
	t.Tweets = append(t.Tweets, twitterScrape(*t)...)
	t.numTweets = len(t.Tweets)
	t.LastScrapeTime = time.Now()
	t.singleComputeHourlySentiment()
}

//DEPRECATED
func (t *ticker) singleComputeHourlySentiment() {
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Print("Error loading sentiment analysis model")
	}
	var total float64
	for _, s := range t.Tweets {
		s.Polarity = sentimentModel.SentimentAnalysis(s.Expression, sentiment.English).Score
		total += float64(s.Polarity)
	}

	t.hourlySentiment = total / float64(t.numTweets)
}*/
