package common

import "encoding/json"

func ToJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ArrayToString(arr []string, separator string) string {
	res := ""
	length := len(arr) - 1
	for i, s := range arr {
		res += s
		if i == length {
			res += ","
		}
	}

	return res
}
