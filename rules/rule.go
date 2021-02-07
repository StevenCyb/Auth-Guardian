package rules

import (
	"auth-guardian/config"
	"auth-guardian/util"
	"net/http"
	"regexp"
)

// Rule define a rule with it's function
type Rule struct {
	// Target attribute
	Method []string
	Path   *regexp.Regexp
}

// FromRuleConfig set attributes from rule config
func (rule *Rule) FromRuleConfig(stringMap config.RuleConfig) {
	rule.Method = stringMap.Method
	if stringMap.Path != "" {
		rule.Path = CreateRegex(stringMap.Path)
	}
}

// TestWhitelist request again this whitelist rule
func (rule *Rule) TestWhitelist(r *http.Request) bool {
	return (len(rule.Method) == 0 || util.StringSliceContains(rule.Method, r.Method)) &&
		(rule.Path == nil || rule.Path.MatchString(r.URL.Path))
}
