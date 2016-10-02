package window

import (
	"time"

	"github.com/idiocracy/mrv/streams"
)

// Window holds window state
type Window struct {
	stream chan streams.Ticker
	size   time.Duration
	state  slidingWindow
}

type slidingWindow map[string]map[time.Time]streams.Ticker

func (sw slidingWindow) addTick(tick streams.Ticker) {
	if ticks, ok := sw[tick.CurrencyPair]; ok {
		ticks[time.Now()] = tick
	} else {
		ticks := make(map[time.Time]streams.Ticker)
		ticks[time.Now()] = tick
		sw[tick.CurrencyPair] = ticks
	}
}

func (sw slidingWindow) assertRetention(retention time.Duration) {
	for _, ticks := range sw {
		for added := range ticks {
			if time.Since(added) > retention {
				delete(ticks, added)
			}
		}
	}
}

// New creates a window aggregator.
// Takes a Ticker channel and a window size as input
func New(stream chan streams.Ticker, size time.Duration) Window {
	state := make(slidingWindow)
	retentionCheck := time.NewTicker(time.Millisecond * 500).C

	go func() {
		for {
			select {
			case <-retentionCheck:
				state.assertRetention(size)
			case tick := <-stream:
				state.addTick(tick)
			}
		}
	}()

	return Window{
		stream: stream,
		size:   size,
		state:  state,
	}
}

// ListCurrencyPairs lists all seen CurrencyPairs
func (w Window) ListCurrencyPairs() (pairs []string) {
	for currencyPair := range w.state {
		pairs = append(pairs, currencyPair)
	}
	return
}

// GetCurrencyPair returns all tickers in current window for a CurrencyPair
func (w Window) GetCurrencyPair(cp string) (t []streams.Ticker) {
	if tickerMap, ok := w.state[cp]; ok {
		for _, ticker := range tickerMap {
			t = append(t, ticker)
		}
	}
	return
}
