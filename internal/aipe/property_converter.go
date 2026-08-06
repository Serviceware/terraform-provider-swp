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
		if v == nil {
			continue
		}

		if b, ok := v.(bool); ok {
			result[k] = fmt.Sprintf("%t", b)
		} else if s, ok := v.(string); ok {
			result[k] = s
		}
		// else: AIPE can return other JSON types (e.g. a nested object like
		// "external-system-reference": {} on manually created objects). We
		// can't represent those as a plain string property, so skip them
		// instead of a panicking type assertion.
	}
	return result
}
