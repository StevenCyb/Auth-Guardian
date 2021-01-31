package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// getConfigFromFile sets the configuration from file
func getConfigFromFile(definition *map[string]map[string]interface{}) {
	// Check config file exists
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return
	}

	// Read config file
	rawYaml, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Panicf("Can't load existing config file %v ", err)
	}

	// Parse config definition
	parsedYaml := map[string]string{}
	err = yaml.Unmarshal(rawYaml, parsedYaml)
	if err != nil {
		log.Panicf("Unmarshaling of existing config file failed %v ", err)
	}

	// Load config file
	for key := range *definition {
		if value, ok := parsedYaml[key]; ok && value != "" {
			(*definition)[key]["file"] = value
		}
	}
}
