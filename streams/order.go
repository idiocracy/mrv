package streams

import "encoding/json"

type Order struct {
	Data json.RawMessage
	Type string
}
