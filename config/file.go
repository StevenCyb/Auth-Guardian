package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// findFirstExistingFile return first existing file in slice
func findFirstExistingFile(fileSlice []string) string {
	for _, path := range fileSlice {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "config.yml"
}

// getConfigFromFile sets the configuration from file
func getConfigFromFile(definition *map[string]map[string]interface{}) error {
	// Find config file to use by prio
	path := findFirstExistingFile([]string{
		"/etc/config/config.yml",
		"/etc/config/config.yaml",
		"config.yml",
		"config.yaml",
	})

	// Check config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Read config file
	rawYaml, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse config definition
	var parsedYaml map[string]interface{}
	err = yaml.Unmarshal(rawYaml, &parsedYaml)
	if err != nil {
		return err
	}

	// Load config file
	for key := range *definition {
		if value, ok := parsedYaml[key]; ok && value != "" {
			(*definition)[key]["file"] = value
		}
	}

	return nil
}
