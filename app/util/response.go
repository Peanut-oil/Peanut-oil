package util

import "time"

func Response(code int, message string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["code"] = code
	result["message"] = message
	result["data"] = data
	result["time"] = time.Now().Unix()
	return result
}
