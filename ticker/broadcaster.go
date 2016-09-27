package ticker

type Broadcaster struct {
	newSubscribers chan chan Ticker
	stories        chan Ticker
}

type TickChan struct {
	C chan Ticker
}

func NewBroadcaster() Broadcaster {
	newSubscribers := make(chan (chan Ticker))
	subscriptions := make([]chan Ticker, 0)
	stories := make(chan Ticker)

	go func() {
		for {
			select {
			case story := <-stories:
				for _, s := range subscriptions {
					s <- story
				}
			case subscription := <-newSubscribers:
				subscriptions = append(subscriptions, subscription)
			}
		}
	}()

	return Broadcaster{
		newSubscribers: newSubscribers,
		stories:        stories,
	}
}

func (b Broadcaster) Listen() TickChan {
	channel := make(chan Ticker)
	b.newSubscribers <- channel
	return TickChan{channel}
}

func (b Broadcaster) Write(tick Ticker) {
	b.stories <- tick
}
