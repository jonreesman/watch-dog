package main

import (
	"log"

	_ "github.com/jonreesman/watch-dog/pb"
)

func main() {
	run()
}

func run() {
	var (
		d DBManager
		b bot
		s Server
	)
	if err := d.initializeManager(); err != nil {
		log.Printf("Failed to initialize DBManager: %v", err)
		return
	}

	if err := b.initBot(d); err != nil {
		log.Printf("Failed to initialize bot: %v", err)
		return
	}
	addTickerChannel := make(chan string)
	deactivateTickerChannel := make(chan int)

	go s.startServer(d, addTickerChannel, deactivateTickerChannel)
	go AddTicker(d, addTickerChannel)
	go DeactivateTicker(d, deactivateTickerChannel)
	b.run(d)
}
