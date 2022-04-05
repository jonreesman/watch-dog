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
		log.Fatalf("Failed to initialize DBManager: %v", err)
	}

	if err := b.initBot(d); err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	go s.startServer(d)
	b.run(d)
}
