package main

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	var d DBManager
	d.initializeManager()
	fmt.Println("Initilialized")
	d.createTickerTable()
	fmt.Println("Ticker Table Created")
	d.createStatementTable()
	fmt.Println("Statement Table Created")
}
