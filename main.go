package main

import (
	"fmt"
	"log"

	"github.com/idiocracy/mrv/ticker"
)

func main() {
	broadcaster := ticker.NewBroadcaster()
	radio := broadcaster.Listen()
	tv := broadcaster.Listen()
	go listenToTheRadio("radio", radio)
	go listenToTheRadio("tv", tv)
	err := ticker.New(broadcaster)
	if err != nil {
		log.Println("[ERROR]", err)
	}

	log.Println("[INFO]", "bye!")
}

// TODO remove this, only for shows
func listenToTheRadio(lol string, r ticker.TickChan) {
	for {
		tick := <-r.C
		fmt.Println(lol, tick)
	}
}
