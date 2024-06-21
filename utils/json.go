package utils

import "encoding/json"

/*
将对象转成字符串
*/
func JsonMarshal(value interface{}) string {
	data, _ := json.Marshal(value)
	return string(data)
}
