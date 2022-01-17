package mdbDriver

import (
	"fmt"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	fmt.Println("RUnning mongotestdb")
	s := statement{
		Expression: "This is a test",
		Subject:    "SPY",
		Source:     "Test",
		TimeStamp:  time.Now().Unix(),
		Polarity:   1,
	}
	mdbPush(s)
}
