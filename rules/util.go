package rules

import (
	"auth-guardian/logging"
	"regexp"
)

// CreateRegex create a regex from given string
func CreateRegex(regString string) *regexp.Regexp {
	reg, err := regexp.Compile(regString)
	if err != nil {
		logging.Fatal(&map[string]string{
			"file":         "util.go",
			"Function":     "CreateRegex",
			"error":        "Regex parsing failed",
			"regex_string": regString,
		})
	}

	return reg
}
