package reception

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

func

func NewPoloniexTicks() error {
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, "wss://api.poloniex.com", nil)

	c.ReceiveTimeout = time.Duration(1 * time.Minute)

	if err != nil {
		log.Println("[DEBUG]", err)
		return fmt.Errorf("Failed to open websocket")
	}

	_, err = c.JoinRealm("realm1", nil)
	if err != nil {
		log.Println("[DEBUG]", err)
		return fmt.Errorf("Unable to join realm")
	}

	log.Println("[INFO]", "Got websocket")

	c.ReceiveDone = make(chan bool)

	err = c.Subscribe("ticker", func(args []interface{}, kwargs map[string]interface{}) {
		fmt.Println("Got tick")
		for _, m := range args {
			fmt.Println(m)
		}

		fmt.Println()
	})

	if err != nil {
		log.Println("[DEBUG]", err)
		return fmt.Errorf("Subscription failed")
	}
	<-c.ReceiveDone
	return nil
}
