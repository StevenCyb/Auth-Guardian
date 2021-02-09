package rules

import (
	"auth-guardian/config"
	"auth-guardian/util"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// Rule define a rule with it's function
type Rule struct {
	// Target attribute
	Method []string
	Path   *regexp.Regexp
	// Condition attribute
	Userinfo          []ValueValidator
	QueryParameter    []ValueValidator
	JSONBodyParameter []ValueValidator
}

// ValueValidator use regex to check value on path
type ValueValidator struct {
	Path      string
	Condition *regexp.Regexp
}

// From set values from path and regex string
func (vv *ValueValidator) From(path string, regexString string) {
	vv.Path = path
	vv.Condition = CreateRegex(regexString)
}

// Match check validator match
func (vv *ValueValidator) Match(fd *util.FlatData) bool {
	data, err := fd.Search(vv.Path)
	if err != nil {
		return false
	}

	switch data.(type) {
	// String related
	case string:
		return vv.Condition.MatchString(data.(string))
	case []string:
		for _, item := range data.([]string) {
			if vv.Condition.MatchString(item) {
				return true
			}
		}
		return false

	// Int related
	case int:
		return vv.Condition.MatchString(fmt.Sprint(data))
	case []int:
		for _, item := range data.([]int) {
			if vv.Condition.MatchString(fmt.Sprint(item)) {
				return true
			}
		}
		return false

	// Float related
	case float32:
		return vv.Condition.MatchString(fmt.Sprint(data))
	case []float32:
		for _, item := range data.([]float32) {
			if vv.Condition.MatchString(fmt.Sprint(item)) {
				return true
			}
		}
		return false

	case float64:
		return vv.Condition.MatchString(fmt.Sprint(data))
	case []float64:
		for _, item := range data.([]float64) {
			if vv.Condition.MatchString(fmt.Sprint(item)) {
				return true
			}
		}
		return false

	// Unknown
	default:
		log.Fatal(&map[string]string{"file": "flatdata.go", "Function": "recursiveBuild", "error": "Unknown type " + fmt.Sprintf("%T", data), "data": fmt.Sprintf("%+v", data)})
	}

	return false
}

// FromRuleConfig set attributes from rule config
func (rule *Rule) FromRuleConfig(ruleConfig config.RuleConfig) {
	rule.Method = ruleConfig.Method

	if ruleConfig.Path != "" {
		rule.Path = CreateRegex(ruleConfig.Path)
	}

	for key, value := range ruleConfig.Userinfo {
		vv := ValueValidator{}
		vv.From(key, value)
		rule.Userinfo = append(rule.Userinfo, vv)
	}

	for key, value := range ruleConfig.QueryParameter {
		vv := ValueValidator{}
		vv.From(key, value)
		rule.QueryParameter = append(rule.Userinfo, vv)
	}

	for key, value := range ruleConfig.JSONBodyParameter {
		vv := ValueValidator{}
		vv.From(key, value)
		rule.JSONBodyParameter = append(rule.Userinfo, vv)
	}
}

// TestWhitelist test request again this whitelist rule
func (rule *Rule) TestWhitelist(r *http.Request) bool {
	return (len(rule.Method) == 0 || util.StringSliceContains(rule.Method, r.Method)) &&
		(rule.Path == nil || rule.Path.MatchString(r.URL.Path))
}

// TestRequired test request again this required rule
func (rule *Rule) TestRequired(r *http.Request, userinfo *util.FlatData, query *util.FlatData, jsonBody *util.FlatData) (bool, bool) {
	if rule.TestWhitelist(r) {
		for _, vv := range rule.Userinfo {
			if !vv.Match(userinfo) {
				return true, false
			}
		}

		for _, vv := range rule.QueryParameter {
			if !vv.Match(query) {
				return true, false
			}
		}

		for _, vv := range rule.JSONBodyParameter {
			if !vv.Match(jsonBody) {
				return true, false
			}
		}

		return true, true
	}
	return false, false
}
