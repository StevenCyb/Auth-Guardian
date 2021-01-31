package config

import "os"

// getConfigFromEnv sets the configuration from env
func getConfigFromEnv(definition *map[string]map[string]interface{}) {
	for key := range *definition {
		if value, ok := os.LookupEnv(key); ok && value != "" {
			(*definition)[key]["env"] = value
		}
	}
}
