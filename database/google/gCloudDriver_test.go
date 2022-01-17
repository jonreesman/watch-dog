package gCloudDriver

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	var d DBManager
	d.InitializeManager()
	fmt.Println("Initilialized")
	d.CreateTickerTable()
	fmt.Println("Ticker Table Created")
	d.CreateStatementTable()
	fmt.Println("Statement Table Created")
}
