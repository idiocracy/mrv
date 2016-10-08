package main

import (
	"fmt"
	"log"
	"time"

	"github.com/idiocracy/mrv/streams"
)

func main() {
	ts, err := streams.NewTickerStream()
	if err != nil {
		log.Println("[ERROR]", err)
	}

	err = streams.NewOrderStream("BTC_XMR")
	if err != nil {
		log.Println("[ERROR]", err)
	}

	// w := window.New(ts.Subscribe().C, time.Minute)

	// TODO following block is only to show the Window API
	for range time.NewTicker(time.Second * 1).C {
		// fmt.Println(w.ListCurrencyPairs())
		// fmt.Println(len(w.GetCurrencyPair("ETH_STEEM")))
	}
	fmt.Println("Started")

	<-ts.RecieveDone
	log.Println("[INFO]", "bye!")
}
