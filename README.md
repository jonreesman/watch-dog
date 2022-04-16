[![CodeFactor](https://www.codefactor.io/repository/github/jonreesman/watch-dog/badge)](https://www.codefactor.io/repository/github/jonreesman/watch-dog)
# watch-dog
The purpose of this project is to scrape twitter, reddit, and various news sources on an hourly basis in order to provide the user with a regular sentiment analysis for their chosen stock/crypto tickers. It's React Native frontend component can currently be found [here](https://github.com/jonreesman/watch-dog-react) which is where the lionshare of my attention presently is focused.

## Front-end
The entire app is capable of adding and removing stocks/crypto on the frontend, displaying the tweets it is currently using as datapoints, and a dual axes graph of the average hourly sentiment analysis overlayed with the price for the same time range. Presently, I have utilized React Native, but will likely shift in the future to try out different frameworks for my own edification.

## Backend
At the beginning of this project, Go had full responsibility for the backend. As time goes on, I've slowly fractured some of the functionalities and implemented them in Python. Go presently is responsible for networking using [Gin](https://github.com/gin-gonic/gin) as well as scraping Twitter using a [Frontend scraper written in Go](https://github.com/n0madic/twitter-scraper) (shoutout to n0madic). The actual Sentiment Analysis and pulling of stock and crypto quotes (the Yahoo Finance API is slowly falling apart and Go packages are casualties) is handled with Python, which communicates with Go via GRPC.

## Database
I explored NoSQL implementations like DynamoDB and MongoDB, but ultimately settled for MySQL. It's tried and true, and I presently don't require the flexibility of NoSQL. As I learn more about Software Engineering however, I find that NoSQL may be a necessity for properly scaling this project should it shift to a centrally run service.

## The Way Forward
- [ ] Reddit Scraping
- [ ] News Scraping

The watch-dog project is my first for Go, so it is a learning curve for me.

