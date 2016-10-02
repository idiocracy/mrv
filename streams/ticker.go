package streams

import (
	"fmt"
	"strconv"

	"gopkg.in/jcelliott/turnpike.v2"
)

const (
	PoloniexWebsocketAddress = "wss://api.poloniex.com"
	PoloniexWebsocketRealm   = "realm1"
	PoloniexWebsocketTopic   = "ticker"
)

type Ticker struct {
	CurrencyPair  string
	Last          float64
	LowestAsk     float64
	HighestBid    float64
	PercentChange float64
	BaseVolume    float64
	QuoteVolume   float64
	IsFrozen      bool
	High          float64
	Low           float64
}

type TickerStream struct {
	B           TickerBroadcaster
	RecieveDone chan bool
}

func NewTickerStream() (ts TickerStream, err error) {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, PoloniexWebsocketAddress, nil)
	if err != nil {
		return ts, fmt.Errorf("Unable to open websocket: %v", err)
	}

	_, err = c.JoinRealm(PoloniexWebsocketRealm, nil)
	if err != nil {
		return ts, fmt.Errorf("Unable to join realm: %v", err)
	}

	c.ReceiveDone = make(chan bool)
	broadcaster := newTickerBroadcaster()

	// TODO Unmarshal the proper ticker event
	err = c.Subscribe(PoloniexWebsocketTopic, handleTicker(broadcaster))
	if err != nil {
		return ts, fmt.Errorf("Subscription to %v failed: %v", PoloniexWebsocketTopic, err)
	}

	return TickerStream{
		B:           broadcaster,
		RecieveDone: c.ReceiveDone,
	}, nil
}

func handleTicker(b TickerBroadcaster) turnpike.EventHandler {
	return func(args []interface{}, kwargs map[string]interface{}) {
		ticker := Ticker{}

		ticker.CurrencyPair = args[0].(string)
		ticker.Last, _ = strconv.ParseFloat(args[1].(string), 64)
		ticker.LowestAsk, _ = strconv.ParseFloat(args[2].(string), 64)
		ticker.HighestBid, _ = strconv.ParseFloat(args[3].(string), 64)
		ticker.PercentChange, _ = strconv.ParseFloat(args[4].(string), 64)
		ticker.BaseVolume, _ = strconv.ParseFloat(args[5].(string), 64)
		ticker.QuoteVolume, _ = strconv.ParseFloat(args[6].(string), 64)

		if args[7].(float64) != 0 {
			ticker.IsFrozen = true
		} else {
			ticker.IsFrozen = false
		}

		ticker.High, _ = strconv.ParseFloat(args[8].(string), 64)
		ticker.Low, _ = strconv.ParseFloat(args[9].(string), 64)

		b.Write(ticker)
	}
}
