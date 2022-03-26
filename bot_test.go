package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTicker(t *testing.T) {
	if os.Getenv("DB_NAME") != "testdb" {
		log.Fatalf("Test database not found in env variables. Exiting.")
	}
	var d DBManager
	var resultSlice tickerSlice
	testTickers := []string{"SPY", "RKLB", "TSLA", "GME",
		"LCID", "NVDA", "TSM", "AMD", "BTC"}
	if err := d.initializeManager(); err != nil {
		log.Printf("Failed to init db")
	}
	for _, tick := range testTickers {
		_, err := AddTicker(d, tick)
		if err != nil {
			log.Printf("Failed to add %s", tick)
		}
	}

	resultSlice = d.returnAllTickers()

	assert.Equal(t, len(resultSlice), len(testTickers), "These should be equal")
	for _, t := range testTickers {
		found := false
		for _, s := range resultSlice {
			if t == s.Name {
				found = true
				break
			}
		}
		if found == false {
			log.Printf("Failed to add %s", t)
		}
	}
}
