package mdbDriver

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type statement struct {
	Expression   string `json:"expression" dynamodbav:"expression" bson:"expression"`
	Subject      string `json:"subject" dynamodbav:"subject" bson:"subject"`
	Source       string `json:"source" dynamodbav:"source" bson:"source"`
	TimeStamp    int64  `json:"timeStamp" dynamodbav:"timestamp" bson:"timestamp"`
	TimeString   string `json:"timeString" dynamodbav:"timeString" bson:"timeString"`
	Polarity     uint8  `json:"polarity" dynamodbav:"polarity" bson:"polarity"`
	timeStampObj time.Time
}

func mdbPush(s statement) {
	fmt.Println("URI IS", os.Getenv("DB_URI"))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("watch-dog").Collection(s.Subject)

	/*docs := []interface{} {
		bson.D{s},
	} */

	/*structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	doc, err := bson.Marshal(s)

	docs := []interface{} {
		bson.D{doc},
	} */

	res, insertErr := collection.InsertOne(ctx, s)

	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)
}
