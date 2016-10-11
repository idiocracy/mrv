package streams

import (
	"encoding/json"
	"strconv"
	"strings"
)

type OrderData struct {
	Rate   float64
	Type   string
	Amount float64
}

func (od *OrderData) UnmarshalJSON(b []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(b, &rawStrings)
	if err != nil {
		return err
	}
	for k, v := range rawStrings {
		if strings.ToLower(k) == "type" {
			od.Type = v
		}
		if strings.ToLower(k) == "rate" {
			od.Rate, _ = strconv.ParseFloat(v, 64)
		}
		if strings.ToLower(k) == "amount" {
			od.Amount, _ = strconv.ParseFloat(v, 64)
		}
	}

	return nil
}
