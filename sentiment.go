package main

import (
	sentiment "github.com/cdipaolo/sentiment"
)

func polarity(tweet string) uint8 {
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}
	result := sentimentModel.SentimentAnalysis(tweet, sentiment.English)
	return result.Score
}
