package config

// RuleConfig define a rule config struct
type RuleConfig struct {
	Type string `json:"type" yaml:"type"`
	// Target attribute
	Method []string `json:"method" yaml:"method"`
	Path   string   `json:"path" yaml:"path"`
	// Condition
	Userinfo          map[string]string `json:"userinfo" yaml:"userinfo"`
	QueryParameter    map[string]string `json:"query-parameter" yaml:"query-parameter"`
	JSONBodyParameter map[string]string `json:"json-body-parameter" yaml:"json-body-parameter"`
}

// HasValidType validate if rule has valid type
func (r *RuleConfig) HasValidType() bool {
	return r.Type == "whitelist" || r.Type == "required" || r.Type == "disallow"
}

// FromMap fill the config rule with data from map
func (r *RuleConfig) FromMap(m map[interface{}]interface{}) {
	for key, value := range m {
		switch key.(string) {
		case "type":
			r.Type = value.(string)
			break
		case "method":
			r.Method = InterfaceToStringSlice(value)
			break
		case "path":
			r.Path = value.(string)
			break
		case "userinfo":
			r.Userinfo = InterfaceToStringMap(value)
			break
		case "query-parameter":
			r.QueryParameter = InterfaceToStringMap(value)
			break
		case "json-body-parameter":
			r.JSONBodyParameter = InterfaceToStringMap(value)
			break
		}
	}
}
