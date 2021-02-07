package config

import (
	"flag"
)

func getConfigFromArguments(definition *map[string]map[string]interface{}) {
	// Register arguments
	for key := range *definition {
		switch (*definition)[key]["type"] {
		case "int":
			(*definition)[key]["arg"] = flag.Int(key, (*definition)[key]["default"].(int), (*definition)[key]["desc"].(string))
			break
		case "bool":
			(*definition)[key]["arg"] = flag.Bool(key, (*definition)[key]["default"].(bool), (*definition)[key]["desc"].(string))
			break
		case "string":
			(*definition)[key]["arg"] = flag.String(key, (*definition)[key]["default"].(string), (*definition)[key]["desc"].(string))
			break
		case "string_array":
			var tmp StringArrayFlag
			(*definition)[key]["arg"] = &tmp
			flag.Var(&tmp, key, (*definition)[key]["desc"].(string))
			break
		case "rule_array":
			var tmp StringArrayFlag
			(*definition)[key]["arg"] = &tmp
			flag.Var(&tmp, key, (*definition)[key]["desc"].(string))
			break
		}
	}

	// Parse arguments
	flag.Parse()

	// Create list of used arguments
	usedArguments := []string{}
	flag.Visit(func(f *flag.Flag) {
		usedArguments = append(usedArguments, f.Name)
	})

	// Remove default values if argument unused
	for key := range *definition {
		if !contains(usedArguments, key) {
			delete((*definition)[key], "arg")
		}
	}
}
