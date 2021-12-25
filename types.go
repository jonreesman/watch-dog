package main

type ticker struct {
	name            string
	numTweets       int
	tweets          []statement
	hourlySentiment float64
}
type statement struct {
	expression string
	polarity   uint8
}
