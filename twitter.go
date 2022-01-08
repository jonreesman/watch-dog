package main

import (
	"context"
	"strings"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func twitterScrape(t ticker) []statement {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	var tweets []statement
	//for tweet := range scraper.SearchTweets(context.Background(), t.name+"-filter:retweets since_time:"+strconv.FormatInt(t.lastScrapeTime.Local().Unix(), 10), 50) {
	for tweet := range scraper.SearchTweets(context.Background(), t.Name+"-filter:retweets within_time:1h", 50) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		if strings.Contains(tweet.Text, t.Name) && (tweet.Timestamp >= (time.Now().Unix() - 3600)) {
			s := statement{
				Expression:   tweet.Text,
				Subject:      t.Name,
				Source:       "Twitter",
				TimeStamp:    tweet.Timestamp,
				timeStampObj: time.Unix(tweet.Timestamp, 0),
				Polarity:     polarity(tweet.Text),
			}
			/*var s statement
			s.Expression = tweet.Text
			s.Subject = t.Name
			s.Source = "Twitter"
			s.TimeStamp = tweet.Timestamp
			s.Polarity = polarity(tweet.Text)*/
			tweets = append(tweets, s)
		}
	}
	return tweets
}
