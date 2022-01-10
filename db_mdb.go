package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
