package utils

import "encoding/json"

func ToJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ToPrettyString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}
