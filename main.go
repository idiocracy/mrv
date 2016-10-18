package main

import (
	"fmt"
	"log"
	"time"

	"github.com/idiocracy/mrv/streams"
	"github.com/idiocracy/mrv/window"
)

func main() {
	ts, err := streams.NewTickerStream()
	if err != nil {
		log.Println("[ERROR]", err)
		panic("Error creating ticker stream, see log")
	}

	os, errOrderStream := streams.NewOrderStream("BTC_XMR")
	if errOrderStream != nil {
		log.Println("[ERROR]", err)
		panic("Error creating order stream, see log")
	}

	orderChan := os.Subscribe().C
	for range time.NewTicker(time.Second * 1).C {
		order := <-orderChan
		fmt.Println(order)
	}

	w := window.New(ts.Subscribe().C, time.Minute)

	// TODO following block is only to show the Window API
	for range time.NewTicker(time.Second * 10).C {
		fmt.Println(w.ListCurrencyPairs())
		fmt.Println(len(w.GetCurrencyPair("ETH_STEEM")))
	}
	fmt.Println("Opened a ticket socket and order socket towards poloniex")

	<-ts.RecieveDone
	log.Println("[INFO]", "bye!")
}
