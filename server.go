package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) startServer(db DBManager, addTicker chan string, deleteTicker chan int) {
	s.addTicker = addTicker
	s.deleteTicker = deleteTicker
	s.d = db

	s.router = gin.Default()

	api := s.router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.HTML(
				http.StatusOK,
				"web/web.html",
				gin.H{"title": "Web"},
			)
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/tickers", s.returnTickersHandler)
		api.POST("/tickers/", s.newTickerHandler)
		api.GET("/tickers/:id/time/:interval", s.returnTickerHandler)
		api.DELETE("/tickers/:id", s.deactivateTickerHandler)

	}

	s.router.Run(":3100")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s Server) newTickerHandler(c *gin.Context) {
	var input ticker
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	s.addTicker <- input.Name
	msg := <-s.addTicker

	c.JSON(http.StatusOK, gin.H{"Id": msg, "Name": input.Name})

}

func (s Server) returnTickersHandler(c *gin.Context) {
	tickers := s.d.returnActiveTickers()
	//Add current prices to tickers
	type payloadItem struct {
		Name            string
		LastScrapeTime  time.Time
		HourlySentiment float64
		Id              int
		Quote           float64
	}
	payload := make([]payloadItem, 0)

	for _, tick := range tickers {
		it := payloadItem{
			Name:            tick.Name,
			LastScrapeTime:  tick.LastScrapeTime,
			HourlySentiment: tick.HourlySentiment,
			Id:              tick.Id,
			Quote:           priceCheck(tick.Name),
		}
		payload = append(payload, it)
	}

	c.JSON(http.StatusOK, payload)
}

func (s Server) returnTickerHandler(c *gin.Context) {
	var (
		id       int
		interval string
		fromTime int64
		t        ticker
		name     string
		period   string
		tick     ticker
		err      error
	)
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id."})
		return
	}
	t, err = s.d.retrieveTickerById(id)
	if err != nil {
		log.Print("Unable to retieve ticker")
	}
	name = t.Name

	interval = c.Param("interval")
	switch interval {
	case "day":
		period = "1d"
	case "week":
		period = "7d"
	case "month":
		period = "30d"
	case "2month":
		period = "60d"
	}
	//fromTime = time.Now().Unix() - int64(period)*3600

	if tick, err = s.d.retrieveTickerById(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve ticker"})
		return
	}

	sentimentHistory := s.d.returnSentimentHistory(id, fromTime)
	quoteHistory, err := http.Get(fmt.Sprintf("http://localhost:8000/ticker/%s/period/%s", name, period))
	body, err := ioutil.ReadAll(quoteHistory.Body)
	log.Print(err)

	var result string
	var jsonResult map[int64]interface{}
	json.Unmarshal([]byte(body), &result)
	json.Unmarshal([]byte(result), &jsonResult)

	statementHistory := s.d.returnAllStatements(id, fromTime)

	c.JSON(http.StatusOK, gin.H{
		"ticker":            tick,
		"quote_history":     jsonResult,
		"sentiment_history": sentimentHistory,
		"statement_history": statementHistory,
	})

}

func (s Server) deactivateTickerHandler(c *gin.Context) {
	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id."})
		return
	}
	s.deleteTicker <- id
	msg := <-s.deleteTicker
	if msg == 400 {
		log.Print("delete failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete ticker"})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": "none"})
	}

}
