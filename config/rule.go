package config

// RuleConfig define a rule config struct
type RuleConfig struct {
	// Target attribute
	Method []string `json:"method" yaml:"method"`
	Path   string   `json:"path" yaml:"path"`
}

// FromMap fill the config rule with data from map
func (r *RuleConfig) FromMap(m map[interface{}]interface{}) {
	for key, value := range m {
		switch key.(string) {
		case "method":
			r.Method = InterfaceToStringSlice(value)
			break
		case "path":
			r.Path = value.(string)
			break
		}
	}
}
