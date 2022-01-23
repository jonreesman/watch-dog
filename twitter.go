package main

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/forPelevin/gomoji"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func twitterScrape(t ticker) []statement {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	scraper.WithReplies(false)
	var tweets []statement
	//for tweet := range scraper.SearchTweets(context.Background(), t.name+"-filter:retweets since_time:"+strconv.FormatInt(t.lastScrapeTime.Local().Unix(), 10), 50) {
	for tweet := range scraper.SearchTweets(context.Background(), t.Name+" within_time:1h", 50) {
		if tweet.Error != nil {
			return tweets
			//panic(tweet.Error)
		}
		tweet.Text = sanitize(tweet.Text)
		if strings.Contains(tweet.Text, t.Name) { //&& (tweet.Timestamp >= (time.Now().Unix() - 3600)) {
			s := statement{
				Expression:   tweet.Text,
				Subject:      t.Name,
				Source:       "Twitter",
				TimeStamp:    tweet.Timestamp,
				TimeString:   time.Unix(tweet.Timestamp, 0).String(),
				timeStampObj: time.Unix(tweet.Timestamp, 0),
				Polarity:     polarity(tweet.Text),
			}
			tweets = append(tweets, s)
		}
	}
	return tweets
}

func sanitize(s string) string {
	emojis := gomoji.FindAll(s)
	regex := regexp.MustCompile("[[:^ascii:]]")

	/* Regex is used here to remove emojis as well.
	 * The accompanying gomoji package is efficient
	 * for finding emojis in the string but not for
	 * removing them.
	 */
	//s = gomoji.RemoveEmojis(s)
	s = regex.ReplaceAllLiteralString(s, "")
	if len(s) == 0 {
		return ""
	}
	for _, emoji := range emojis {
		s += " " + emoji.Slug + " "
	}
	s = strings.ReplaceAll(s, "\"", "'")
	s = strings.ReplaceAll(s, ";", "semi-colon")
	return s
}
