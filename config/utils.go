package config

import "strconv"

// Check if string array contains string
func contains(slice []string, sa string) bool {
	for _, sb := range slice {
		if sa == sb {
			return true
		}
	}
	return false
}

// getMostlyPrioriesConfigKey get key for the priories config value
func getMostlyPrioriesConfigKey(option map[string]interface{}) interface{} {
	if value, ok := option["arg"]; ok {
		if option["type"] == "int" {
			return *(value.(*int))
		} else if option["type"] == "bool" {
			return *(value.(*bool))
		}
		return *(value.(*string))
	} else if value, ok := option["file"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i
		} else if option["type"] == "bool" {
			return value.(string) == "true"
		}
		return value
	} else if value, ok := option["env"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i
		} else if option["type"] == "bool" {
			return value.(string) == "true"
		}
		return value
	}
	return option["default"]
}
