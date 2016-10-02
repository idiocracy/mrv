package main

import (
	"fmt"
	"log"

	"github.com/idiocracy/mrv/streams"
)

func main() {
	ts, err := streams.NewTickerStream()
	if err != nil {
		log.Println("[ERROR]", err)
	}

	radio := ts.Subscribe()

	go listenToTheRadio(radio)

	<-ts.RecieveDone
	log.Println("[INFO]", "bye!")
}

// TODO remove this, only for shows
func listenToTheRadio(r streams.TickChan) {
	for {
		tick := <-r.C
		fmt.Println(tick)
	}
}
