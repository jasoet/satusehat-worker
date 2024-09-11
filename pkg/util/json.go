package util

import (
	"encoding/json"
)

func MarshalToJson(o any) *json.RawMessage {
	jsonData, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	rawJson := json.RawMessage(jsonData)
	return &rawJson
}
