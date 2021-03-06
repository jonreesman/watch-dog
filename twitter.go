package main

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/forPelevin/gomoji"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func twitterScrape(t ticker) []statement {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	scraper.WithReplies(false)
	var tweets []statement
	for tweet := range scraper.SearchTweets(context.Background(), t.Name+" within_time:1h", 100) {
		if tweet.Error != nil {
			return tweets
		}
		tweet.Text = sanitize(tweet.Text)
		if strings.Contains(tweet.Text, t.Name) {
			if tweet.Timestamp < t.LastScrapeTime.Unix() {
				break
			}
			id, err := strconv.ParseUint(tweet.ID, 10, 64)
			if err != nil {
				log.Printf("twitterScrape(): Error extracting tweet id")
			}
			s := statement{
				Expression:   tweet.Text,
				subject:      t.Name,
				source:       "Twitter",
				TimeStamp:    tweet.Timestamp,
				Polarity:     0,
				URLs:         tweet.URLs,
				PermanentURL: tweet.PermanentURL,
				ID:           id,
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
