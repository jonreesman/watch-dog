package main

import (
	"testing"
	"time"
)

func Test_grabQuotes(t *testing.T) {
	type args struct {
		d             DBManager
		quoteInterval time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//grabQuotes(tt.args.d, tt.args.quoteInterval)
		})
	}
}
