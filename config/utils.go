package config

import (
	"encoding/json"
	"log"
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
	// Parse arguments
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

		} else if option["type"] == "rule_array" {
			rules := []RuleConfig{}
			for _, item := range *(value.(*StringArrayFlag)) {
				rule := RuleConfig{}
				err := json.Unmarshal([]byte(item), &rule)
				if err != nil {
					log.Panic(map[string]string{
						"log_level": "panic",
						"error":     "Rule unmarshal failed",
						"for_rule":  item,
						"details":   err.Error(),
					})
				}
				rules = append(rules, rule)
			}
			return rules

		} else if option["type"] == "string" {
			return *(value.(*string))
		}

		// Parse config file
	} else if value, ok := option["file"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i

		} else if option["type"] == "bool" {
			return value.(string) == "true"

		} else if option["type"] == "string_array" {
			return InterfaceToStringSlice(value)

		} else if option["type"] == "rule_array" {
			rules := []RuleConfig{}
			for _, v := range value.([]interface{}) {
				rule := RuleConfig{}
				rule.FromMap(v.(map[interface{}]interface{}))
				rules = append(rules, rule)
			}
			return rules

		} else if option["type"] == "string" {
			return value.([]map[string]string)
		}

		// Parse environment
	} else if value, ok := option["env"]; ok {
		if option["type"] == "int" {
			i, _ := strconv.Atoi(value.(string))
			return i

		} else if option["type"] == "bool" {
			return value.(string) == "true"

		} else if option["type"] == "string_array" {
			return strings.Split(value.(string), ",")

		} else if option["type"] == "rule_array" {
			rules := []RuleConfig{}
			err := json.Unmarshal([]byte(value.(string)), &rules)
			if err != nil {
				log.Panic(map[string]string{
					"log_level": "panic",
					"error":     "Rules unmarshal failed",
					"for_rule":  value.(string),
					"details":   err.Error(),
				})
			}
			return rules

		} else if option["type"] == "string" {
			return value.(string)
		}
	}

	// Send back default
	if option["type"] == "string_array" {
		tmp := []string{}
		for _, item := range *(option["default"].(*StringArrayFlag)) {
			tmp = append(tmp, item)
		}
		return tmp
	}

	if option["type"] == "rule_array" {
		return []RuleConfig{}
	}

	return option["default"]
}

// StringArrayFlag is a string array for flag
type StringArrayFlag []string

// String return array as string
func (i *StringArrayFlag) String() string {
	return strings.Join(*i, " ")
}

// Set a item to array
func (i *StringArrayFlag) Set(value string) error {
	var jData1 map[string]interface{}
	if err := json.Unmarshal([]byte(value), &jData1); err == nil {
		*i = append(*i, value)
		return nil
	}
	var jData2 []string
	if err := json.Unmarshal([]byte(value), &jData2); err == nil {
		for _, item := range jData1 {
			stringData, _ := json.Marshal(item)
			*i = append(*i, string(stringData))
		}
		return nil
	}
	if strings.Contains(value, ",") {
		for _, v := range strings.Split(value, ",") {
			*i = append(*i, v)
		}
		return nil
	}

	*i = append(*i, value)
	return nil
}

// InterfaceToStringSlice convert a interface to string slice
func InterfaceToStringSlice(i interface{}) []string {
	stringSlice := []string{}
	for _, v := range i.([]interface{}) {
		stringSlice = append(stringSlice, v.(string))
	}
	return stringSlice
}

// InterfaceToStringMap convert a interface to string map
func InterfaceToStringMap(i interface{}) map[string]string {
	stringMap := map[string]string{}
	for key, value := range i.(map[interface{}]interface{}) {
		stringMap[key.(string)] = value.(string)
	}
	return stringMap
}
