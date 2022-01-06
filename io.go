package main

import (
	"encoding/json"
	"log"
	"io"
)

func (t ticker) dump() {
	p, err := json.Marshall(t.tweets)
	if err != nil {
		log.Panicf("Error dumping %s\n", err)
	}
	file, err := io.WriteFile(t.name + ".json",p,0644)
	if err != nil {
		log.Panicf("faild reading data from file: %s", err)
	}
	defer file.Close()

}
