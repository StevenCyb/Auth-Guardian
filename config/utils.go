package config

import (
	"fmt"
	"strconv"
	"strings"
)

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
		} else if option["type"] == "string_array" {
			tmp := []string{}
			for _, item := range *(value.(*StringArrayFlag)) {
				tmp = append(tmp, item)
			}
			return tmp
		}
		return *(value.(*string))
	} else if value, ok := option["file"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i
		} else if option["type"] == "bool" {
			return value.(string) == "true"
		} else if option["type"] == "string_array" {
			valueSlice := value.([]interface{})
			stringSlice := []string{}
			for _, v := range valueSlice {
				stringSlice = append(stringSlice, v.(string))
			}
			return stringSlice
		}
		return value
	} else if value, ok := option["env"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i
		} else if option["type"] == "bool" {
			return value.(string) == "true"
		} else if option["type"] == "string_array" {
			return strings.Split(value.(string), ",")
		}
		return value
	}

	// Send back default
	if option["type"] == "string_array" {
		fmt.Println(option["default"].(*StringArrayFlag))
		tmp := []string{}
		for _, item := range *(option["default"].(*StringArrayFlag)) {
			tmp = append(tmp, item)
		}
		return tmp
	}
	return option["default"]
}

// StringArrayFlag is a string array for flag
type StringArrayFlag []string

// String return array as string
func (i *StringArrayFlag) String() string {
	return ""
}

// Set a item to array
func (i *StringArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
