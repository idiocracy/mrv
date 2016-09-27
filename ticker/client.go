package ticker

import (
	"fmt"

	"gopkg.in/jcelliott/turnpike.v2"
)

const (
	POLONIEX_WEBSOCKET_ADDRESS = "wss://api.poloniex.com"
	POLONIEX_WEBSOCKET_REALM   = "realm1"
	POLONIEX_WEBSOCKET_TICKER  = "ticker"
)

// TODO return a TickerConnection struct or the like
func New(broadcaster Broadcaster) error {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, POLONIEX_WEBSOCKET_ADDRESS, nil)
	if err != nil {
		return fmt.Errorf("Unable to open websocket: %v", err)
	}

	_, err = c.JoinRealm(POLONIEX_WEBSOCKET_REALM, nil)
	if err != nil {
		return fmt.Errorf("Unable to join realm: %v", err)
	}

	c.ReceiveDone = make(chan bool)

	// TODO Unmarshal the proper ticker event
	err = c.Subscribe(POLONIEX_WEBSOCKET_TICKER, func(args []interface{}, kwargs map[string]interface{}) {
		broadcaster.Write(Ticker{args[0].(string)})
	})
	if err != nil {
		return fmt.Errorf("Subscription to %v failed: %v", POLONIEX_WEBSOCKET_TICKER, err)
	}

	<-c.ReceiveDone
	return nil
}
