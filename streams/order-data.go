package streams

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type orderDataRaw struct {
	Rate    string `json:"rate"`
	Type    string `json:"type"`
	Amount  string `json:"amount"`
	TradeID string `json:"tradeID"`
	Date    string `json:"date"`
	Total   string `json:"total"`
}

type OrderData struct {
	Rate    float64
	Type    string
	Amount  float64
	TradeID string
	Date    string
	Total   float64
}

func (od *OrderData) UnmarshalJSON(b []byte) error {
	var raw orderDataRaw

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	od.Type = raw.Type

	if raw.TradeID != "" {
		od.TradeID = raw.TradeID
	}

	if raw.Date != "" {
		od.Date = raw.Date
	}

	if raw.Rate != "" {
		od.Rate, err = strconv.ParseFloat(raw.Rate, 64)
		if err != nil {
			return fmt.Errorf("Could not parse orderDataRaw.Rate as float64: %v", err)
		}
	}

	if raw.Amount != "" {
		od.Amount, err = strconv.ParseFloat(raw.Amount, 64)
		if err != nil {
			return fmt.Errorf("Could not parse orderDataRaw.Amount as float64: %v", err)
		}
	}

	if raw.Total != "" {
		od.Total, err = strconv.ParseFloat(raw.Total, 64)
		if err != nil {
			return fmt.Errorf("Could not parse orderDataRaw.Total as float64: %v", err)
		}
	}

	return nil
}
