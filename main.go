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
		return
	}

	if err := b.initBot(d); err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
		return
	}

	go s.startServer(d)
	b.run(d)
}
