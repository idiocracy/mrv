package streams

type OrderBroadcaster struct {
	newSubscribers chan chan Order
	stories        chan Order
}

type OrderChan struct {
	C chan Order
}

func newOrderBroadcaster() OrderBroadcaster {
	newSubscribers := make(chan (chan Order))
	var subscriptions []chan Order
	stories := make(chan Order)

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

	return OrderBroadcaster{
		newSubscribers: newSubscribers,
		stories:        stories,
	}
}

func (b OrderBroadcaster) Listen() OrderChan {
	channel := make(chan Order)
	b.newSubscribers <- channel
	return OrderChan{channel}
}

func (b OrderBroadcaster) Write(order Order) {
	b.stories <- order
}
