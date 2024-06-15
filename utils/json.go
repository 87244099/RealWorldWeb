package utils

import "encoding/json"

func JsonMarshal(value interface{}) string {
	data, _ := json.Marshal(value)
	return string(data)
}
