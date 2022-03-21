[![CodeFactor](https://www.codefactor.io/repository/github/jonreesman/watch-dog/badge)](https://www.codefactor.io/repository/github/jonreesman/watch-dog)
# watch-dog
The purpose of this project is to scrape twitter, reddit, and various news sources on an hourly basis in order to provide the user with a regular sentiment analysis for their chosen stock/crypto tickers. It will feature a frontend, built with React, that will chart stock prices vs their hourly sentiment on Twitter. As new stocks/cryptos are added, they will be cached in the database to reduce the API calls to Twitter. As the program detects headlines that are considered popular for a given stock/crypto (eg, elevated volume of retweeting the article), it will incorporate these headlines into the graph as well.

While I do not believe this project will give the user a leg up in terms of day-trading, I do think it possesses value as a market research tool. Some stocks and crypto have a low volume of tweets, making it difficult to glean any time-sensitive, useful data. But being able to graph sentiment over time and correlate it with the price has been a large desire of mine, because I do think there is value there. In it's most simple form, it may be able to help the user detect potential pump and dumps, which is why I named it watch-dog.

This is a very early build, that presently can scrape twitter and provide a sentiment from its extracted tweets and print it to stdout. I have experimented with multiple different databases (DynamoDB, MongoDB) but have ultimately settled on MySQL via Google Cloud.

The watch-dog project is my first for Go, so it is a learning curve for me.

## To Do:
- [ ] Implement a Web UI (REACT)
- [ ] Incorporate relevant news scraping
