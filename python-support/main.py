import string
from unittest import result
from fastapi import FastAPI
from pydantic import BaseModel
import uvicorn
import json
from textblob import TextBlob

import pandas as pd
import yfinance as yf


app = FastAPI()

@app.get("/")
async def read_root():
    return {"Hello": "World"}

class Data(BaseModel):
    tweet: str

@app.post("/tweet")
async def sentiment_analysis(data: Data):
    print(data.tweet)
    blob = TextBlob(data.tweet)
    print(blob.sentiment.polarity)
    return blob.sentiment.polarity

#1d->60d

@app.get("/ticker/{ticker_name}/period/{period}")
def market_history(ticker_name: str, period: str):
    data = yf.download(tickers=ticker_name, period=period, interval='5m')
    data = data.drop(['High','Low','Close','Adj Close', 'Volume'], axis=1)
    print(data.columns)
    print(data)
    #return data.to_json()
    result = data.to_json(orient="index")
    print(result)
    return result

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)