package streams

import (
	"encoding/json"
	"fmt"

	"gopkg.in/jcelliott/turnpike.v2"
)

// TickerStream holds state for a Poloniex 'ticker' firehose
type OrderStream struct {
	broadcaster TickerBroadcaster
	RecieveDone chan bool
}

// NewTickerStream connects to the Poloniex 'ticker' firehose.
func NewOrderStream(market string) (err error) {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, poloniexWebsocketAddress, nil)
	if err != nil {
		return fmt.Errorf("Unable to open websocket: %v", err)
	}

	_, err = c.JoinRealm(poloniexWebsocketRealm, nil)
	if err != nil {
		return fmt.Errorf("Unable to join realm: %v", err)
	}

	c.ReceiveDone = make(chan bool)

	// TODO Unmarshal the proper ticker event
	err = c.Subscribe(market, handleOrder())
	if err != nil {
		return fmt.Errorf("Subscription to %v failed: %v", poloniexWebsocketTopic, err)
	}
	return nil
}

func handleOrder() turnpike.EventHandler {
	return func(args []interface{}, kwargs map[string]interface{}) {
		order := Order{}
		for _, v := range args {
			str, err := json.Marshal(v)
			if err != nil {
				fmt.Println("Error encoding JSON")
				return
			}
			json.Unmarshal([]byte(str), &order)
			fmt.Println(order)
		}
	}
}
