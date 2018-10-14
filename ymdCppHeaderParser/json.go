package ymdCppHeaderParser

import "encoding/json"

func marshalJson(v interface{}) string {
	if v == nil {
		return `<nil>`
	}
	data, _ := json.Marshal(v)
	return string(data)
}
