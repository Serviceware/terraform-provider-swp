package aipe

import "fmt"

func convertPropertiesFromString(properties map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range properties {
		if v == "false" {
			result[k] = false
		} else if v == "true" {
			result[k] = true
		} else {
			result[k] = v
		}
	}
	return result
}

func convertPropertiesToString(properties map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range properties {
		if b, ok := v.(bool); ok {
			result[k] = fmt.Sprintf("%t", b)
		} else {
			result[k] = v.(string)
		}
	}
	return result
}
