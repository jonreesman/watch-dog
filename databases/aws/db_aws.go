package awsDriver

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"log"
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

func awsPush(s statement) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	svc := dynamodb.New(sess)

	item := s

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)

	}
	tableName := "Tickers"
	fmt.Println(av)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}