package window

import (
	"time"

	"github.com/idiocracy/mrv/streams"
)

// Window holds window state
type WindowOrders struct {
	stream       chan streams.Order
	size         time.Duration
	state        slidingWindowOrders
	currencyPair string
}

// New creates a window aggregator.
// Takes a Ticker channel and a window size as input
func NewOrderWindow(stream chan streams.Order, size time.Duration, currencyPair string) WindowOrders {
	state := make(slidingWindowOrders)
	retentionCheck := time.NewTicker(time.Millisecond * 500).C

	go func() {
		for {
			select {
			case <-retentionCheck:
				state.assertRetention(size)
			case order := <-stream:
				state.addOrder(order)
			}
		}
	}()

	return WindowOrders{
		stream:       stream,
		size:         size,
		state:        state,
		currencyPair: currencyPair,
	}
}

// GetCurrencyPair returns all tickers in current window for a CurrencyPair
func (wo WindowOrders) GetOrdersFor(orderType string) (o []streams.Order) {
	if orderMap, ok := wo.state[orderType]; ok {
		for _, order := range orderMap {
			o = append(o, order)
		}
	}
	return
}

type slidingWindowOrders map[string]map[time.Time]streams.Order

func (sw slidingWindowOrders) addOrder(order streams.Order) {
	if orders, ok := sw[order.Data.Type]; ok {
		orders[time.Now()] = order
	} else {
		orders := make(map[time.Time]streams.Order)
		orders[time.Now()] = order
		sw[order.Data.Type] = orders
	}
}

func (sw slidingWindowOrders) assertRetention(retention time.Duration) {
	for _, orders := range sw {
		for added := range orders {
			if time.Since(added) > retention {
				delete(orders, added)
			}
		}
	}
}
