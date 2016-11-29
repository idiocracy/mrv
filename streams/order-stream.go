package streams

import (
	"encoding/json"
	"fmt"

	"gopkg.in/jcelliott/turnpike.v2"
)

// TickerStream holds state for a Poloniex 'ticker' firehose
type OrderStream struct {
	broadcaster OrderBroadcaster
	RecieveDone chan bool
}

// NewTickerStream connects to the Poloniex 'ticker' firehose.
func NewOrderStream(market string) (os OrderStream, err error) {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, poloniexWebsocketAddress, nil)
	if err != nil {
		return os, fmt.Errorf("Unable to open websocket: %v", err)
	}

	_, err = c.JoinRealm(poloniexWebsocketRealm, nil)
	if err != nil {
		return os, fmt.Errorf("Unable to join realm: %v", err)
	}

	c.ReceiveDone = make(chan bool)
	broadcaster := newOrderBroadcaster()

	// TODO Unmarshal the proper ticker event
	err = c.Subscribe(market, handleOrder(broadcaster))
	if err != nil {
		return os, fmt.Errorf("Subscription to %v failed: %v", poloniexWebsocketTopic, err)
	}
	return OrderStream{
		broadcaster: broadcaster,
		RecieveDone: c.ReceiveDone,
	}, nil

}

func (os OrderStream) Subscribe() OrderChan {
	return os.broadcaster.Listen()
}

func handleOrder(b OrderBroadcaster) turnpike.EventHandler {
	return func(args []interface{}, kwargs map[string]interface{}) {
		order := Order{}
		for _, v := range args {
			str, err := json.Marshal(v)
			if err != nil {
				fmt.Println("Error encoding JSON")
				return
			}
			errUnmarshal := json.Unmarshal([]byte(str), &order)
			if errUnmarshal != nil {
				fmt.Println(errUnmarshal)
			}
			b.Write(order)
		}
	}
}
