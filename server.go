package main

import (
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

	/*var tickerPackage []tickerPayLoad
	for _, tick := range *s.t {
		tickerPackage = append(tickerPackage, tickerPayLoad{Id: tick.id, Name: tick.Name})
	}*/
	c.JSON(http.StatusOK, tickers)
}

func (s Server) returnTickerHandler(c *gin.Context) {
	var (
		id       int
		interval string
		fromTime int64
		hours    int
		tick     ticker
		err      error
	)
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id."})
		return
	}

	interval = c.Param("interval")
	switch interval {
	case "day":
		hours = 24
	case "week":
		hours = 168
	case "month":
		hours = 730
	case "3month":
		hours = 2190
	case "6month":
		hours = 4380
	case "year":
		hours = 8760
	}
	fromTime = time.Now().Unix() - int64(hours)*3600

	if tick, err = s.d.retrieveTickerById(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve ticker"})
		return
	}

	sentimentHistory := s.d.returnSentimentHistory(id, fromTime)
	quoteHistory := s.d.returnQuoteHistory(id, fromTime)
	statementHistory := s.d.returnAllStatements(id, fromTime)

	c.JSON(http.StatusOK, gin.H{
		"ticker":            tick,
		"quote_history":     quoteHistory,
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
