package streams

type TickerBroadcaster struct {
	newSubscribers chan chan Ticker
	stories        chan Ticker
}

type TickChan struct {
	C chan Ticker
}

func newTickerBroadcaster() TickerBroadcaster {
	newSubscribers := make(chan (chan Ticker))
	var subscriptions []chan Ticker
	stories := make(chan Ticker)

	go func() {
		for {
			select {
			case subscription := <-newSubscribers:
				subscriptions = append(subscriptions, subscription)
			case story := <-stories:
				for _, s := range subscriptions {
					s <- story
				}
			}
		}
	}()

	return TickerBroadcaster{
		newSubscribers: newSubscribers,
		stories:        stories,
	}
}

func (b TickerBroadcaster) Listen() TickChan {
	channel := make(chan Ticker)
	b.newSubscribers <- channel
	return TickChan{channel}
}

func (b TickerBroadcaster) Write(tick Ticker) {
	b.stories <- tick
}
