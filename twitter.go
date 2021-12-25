package main

import (
	"context"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"strings"
)

//twitter.go
//context
//twitterscraper
//strings
func twitterScrape(n string) []statement {
	scraper := twitterscraper.New()
	var tweets []statement
	for tweet := range scraper.SearchTweets(context.Background(), n+"- filter:retweets", 50) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		if strings.Contains(tweet.Text, n) == true {
			var s statement
			s.expression = tweet.Text
			s.polarity = polarity(tweet.Text)
			tweets = append(tweets, s)
		}
	}
	return tweets
}
