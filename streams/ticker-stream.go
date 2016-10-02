package streams

import (
	"fmt"
	"strconv"

	"gopkg.in/jcelliott/turnpike.v2"
)

const (
	poloniexWebsocketAddress = "wss://api.poloniex.com"
	poloniexWebsocketRealm   = "realm1"
	poloniexWebsocketTopic   = "ticker"
)

// TickerStream holds state for a Poloniex 'ticker' firehose
type TickerStream struct {
	broadcaster TickerBroadcaster
	RecieveDone chan bool
}

// NewTickerStream connects to the Poloniex 'ticker' firehose.
func NewTickerStream() (ts TickerStream, err error) {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, poloniexWebsocketAddress, nil)
	if err != nil {
		return ts, fmt.Errorf("Unable to open websocket: %v", err)
	}

	_, err = c.JoinRealm(poloniexWebsocketRealm, nil)
	if err != nil {
		return ts, fmt.Errorf("Unable to join realm: %v", err)
	}

	c.ReceiveDone = make(chan bool)
	broadcaster := newTickerBroadcaster()

	// TODO Unmarshal the proper ticker event
	err = c.Subscribe(poloniexWebsocketTopic, handleTicker(broadcaster))
	if err != nil {
		return ts, fmt.Errorf("Subscription to %v failed: %v", poloniexWebsocketTopic, err)
	}

	return TickerStream{
		broadcaster: broadcaster,
		RecieveDone: c.ReceiveDone,
	}, nil
}

// Subscribe to the Ticker stream from Poloniex.
// Method returns a TickChan containing a channel with all new messages on the Poloniex 'ticker' stream.
func (ts TickerStream) Subscribe() TickChan {
	return ts.broadcaster.Listen()
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
