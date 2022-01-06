package main

import (
	"context"
	"strings"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func twitterScrape(n string) []statement {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	var tweets []statement
	for tweet := range scraper.SearchTweets(context.Background(), n+"-filter:retweets", 50) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		if strings.Contains(tweet.Text, n) == true {
			var s statement
			s.expression = tweet.Text
			s.subject = n
			s.source = "Twitter"
			s.timeStamp = tweet.Timestamp
			s.polarity = polarity(tweet.Text)
			tweets = append(tweets, s)
		}
	}
	return tweets
}
