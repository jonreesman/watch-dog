# watch-dog
The purpose of this project is to scrape twitter, reddit, and various news sources on an hourly basis in order to provide the user with a regular sentiment analysis for their chosen stock/crypto tickers.

This is a very early build, that presently can scrape twitter and provide a sentiment from its extracted tweets and print it to stdout. It does have the capability to push to an AWS DynamoDB, however my additions to my DYnamoDB table are throttled, so I am in the progress of swithcing to a local MongoDB table. I have a lengthy to-do list, and the actual structure of the project is a WIP as I am still learning Go best practices.

The watch-dog project is my first for Go, so it is a learning project for me.

The current to do list:
1. MongoDB Integration
2. Incorporate Reddit scraping (I have yet to find an open source project I would like to use)
3. Incorporate relevant news scraping
4. Refactor in accordance with Go best practice 
5. Collect a large data set of tweets, reddit comments, news sources to use to more applicably train the present sentiment analysis model.
6. Implement a Web UI or form of notification system

Presently, the porject is structured as so:

main.go: Creates a Bot instance and spins it up utilizing a Bot reciever function.

bot.go: Implements the business logic of the Bot type.

types.go: Defines the different objects used throughout the project.
  ticker - Stores the name of the stock/crypto ticker, the time of the last web scraping, all the tweets for that ticker, and the hourly sentiment.
  statement - Stores the data for each tweet (defined here as "expression" for flexibility when reddit and news scraping is implemented). It will also be flexible and store reddit comments and news articles. It also stores the polarity of the expression.
  bot - Stores our slice of tickers, as well as the time interval in which scraping will be conducted (for flexibility while debugging, you can reduce the interval from hourly to your own short interval).
  
ticker.go: Handles the business logic of the ticker type.
  computeHourlySentiment() float64
      Returns the average sentiment (defined as polarity) of all tweets currently stored for the specific ticker.
      
  pushToDb()
      Presently, this is a DynamoDb specific function. Makes calls to db.go to push the tweets of the specific ticker to your DynamoDb table.
  
  printTicker()
      Primarily a debug function, but it prints out the name of the ticker, the number of tweets presently stored, the time of the last scrape, and all the tweets it currently has.
  
  importTickers() []ticker
      Basic IO function that takes your chosen tickers from the "tickers.txt" file. Returns a slice of ticker objects, one for each stock/crypto ticker.
      
  scrapeAll(t *[]ticker)
      Handles the slice of Tickers (imported by importTickers()), and calls the scrape receiver function for each ticker.
  
  scrape()
      Receiver function for individual ticker objects. Makes calls to twitter.go to scrape and store tweets in the tickers object.
      
twitter.go: Handles the business logic of calling the twitterscraper package (github.com/n0madic/twitter-scraper).
  twitterScraper(n string) []statement
      Takes the name of the ticker (eg: SPY, TSLA, VOO, etc), and makes calls to the twitterscraper package to scrape the tweets. It then builds the statement object and returns it to the scrape() ticker receiver function.
      
sentiment.go: Utilizes the sentiment package (github.com/cdipaolo/sentiment) to run sentiment analysis on the tweet it is passed.
  polarity(tweet string) uint8
      Takes the tweet as a string, runs sentiment analysis, and returns the polarity (0: negative, 1: positive) representing the tweets sentiment on the ticker.
