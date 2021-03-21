package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// getConfigFromFile sets the configuration from file
func getConfigFromFile(definition *map[string]map[string]interface{}) error {
	// Check config file exists
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return nil
	}

	// Read config file
	rawYaml, err := ioutil.ReadFile("config.yml")
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
