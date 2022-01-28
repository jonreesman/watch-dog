package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	d      *DBManager
	t      *[]ticker
	router *gin.Engine
}

func (s *Server) startServer(db *DBManager, ti *[]ticker) {
	port := os.Getenv("PORT")
	s.d = db
	s.t = ti
	if port == "" {
		port = "3050"
	}
	s.router = gin.Default()
	//s.router.Use(static.Serve("/", static.LocalFile("./views", true)))
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
		api.GET("/tickers", s.returnTickers)
		api.POST("/tickers/", s.newTicker)
		api.GET("/tickers/:id/time/:interval", s.returnTicker)

	}
	//api.GET("/tickers", tickerHandler)
	s.router.Run(":3100")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s Server) newTicker(c *gin.Context) {
	var input ticker
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	newTickerObject, err := addTicker(input.Name, *s.d)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	*s.t = append(*s.t, newTickerObject)
	c.JSON(http.StatusOK, gin.H{"Id": newTickerObject.id, "Name": newTickerObject.Name})

}

func (s Server) returnTickers(c *gin.Context) {
	//tickers := s.d.retrieveTickers()

	type tickerPayLoad struct {
		Name string
		Id   int
	}

	var tickerPackage []tickerPayLoad
	for _, tick := range *s.t {
		tickerPackage = append(tickerPackage, tickerPayLoad{Id: tick.id, Name: tick.Name})
	}
	c.JSON(http.StatusOK, tickerPackage)
}

/*	24 hours 	24 * time.Hour
 *	1 week		168 * time.Hour
 *	1 month		730 * time.Hour
 *  3 months	2190 * time.Hour
 *	6 months	4380 * time.Hour
 *	1 year		8760 * time.Hour
 */
func (s Server) returnTicker(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"ticker":            tick,
		"quote_history":     quoteHistory,
		"sentiment_history": sentimentHistory,
	})

}

func (s Server) getTickerSentimentTimeline(c *gin.Context) {

}