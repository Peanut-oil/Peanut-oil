package serialize

func Response(code int, message string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["code"] = code
	result["message"] = message
	result["data"] = data
	return result
}
