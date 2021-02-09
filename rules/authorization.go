package rules

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/util"
	"fmt"
	"io/ioutil"
	"net/http"
)

var requiredRules []Rule
var disallowRules []Rule
var ruleUseBody bool
var ruleUseQuery bool

// InitializeAuthorizationMiddleware initialize the whitelist middleware
func InitializeAuthorizationMiddleware() {
	logging.Debug(&map[string]string{"file": "authorization.go", "Function": "InitializeAuthorizationMiddleware", "event": "Initialize authorization middleware"})

	ruleUseBody = false
	ruleUseQuery = false

	// Initialize required rules
	for _, rule := range config.Rules {
		ruleName := ""

		ruleStruct := Rule{}
		ruleStruct.FromRuleConfig(rule)
		if rule.Type == "required" {
			ruleName = "Required"
			requiredRules = append(requiredRules, ruleStruct)
		} else if rule.Type == "disallow" {
			ruleName = "Disallow"
			disallowRules = append(requiredRules, ruleStruct)
		}

		if len(rule.JSONBodyParameter) > 0 {
			ruleUseBody = true
		}

		if len(rule.QueryParameter) > 0 {
			ruleUseQuery = true
		}

		logging.Info(&map[string]string{
			"event":                    ruleName + " rule added",
			"rule_method":              fmt.Sprintf("%v", rule.Method),
			"rule_path":                rule.Path,
			"rule_userinfo":            fmt.Sprintf("%v", rule.Userinfo),
			"rule_query_parameter":     fmt.Sprintf("%v", rule.QueryParameter),
			"rule_json_body_parameter": fmt.Sprintf("%v", rule.JSONBodyParameter),
		})
	}
}

// AuthorizationRuleMiddleware check request are authorized or unauthorized through rules
func AuthorizationRuleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body
		bodyString := ""
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			bodyString = string(body)
		}

		logging.Debug(&map[string]string{
			"file":      "authorization.go",
			"Function":  "AuthorizationRuleMiddleware",
			"event":     "Handle request",
			"method":    r.Method,
			"req_path":  r.URL.Path,
			"req_query": r.URL.RawQuery,
			"req_body":  bodyString,
		})

		bodyFD := &util.FlatData{}
		queryFD := &util.FlatData{}
		userinfo := &util.FlatData{}
		skip := len(requiredRules) == 0 && len(disallowRules) == 0

		if !skip {
			if ruleUseBody {
				bodyMap, _ := util.JSONToMap(string(body))
				bodyFD = util.NewFlatData()
				bodyFD.BuildFrom(bodyMap)
			}

			if ruleUseQuery {
				queryMap := map[string]interface{}{}
				for key, value := range r.URL.Query() {
					queryMap[key] = value
				}
				queryFD = util.NewFlatData()
				queryFD.BuildFrom(queryMap)
			}

			if sessionID := r.Context().Value(util.ContextKey("session_id")); sessionID != nil {
				session := util.SessionMap[sessionID.(string)]

				userinfoI, err := session.Get("userinfo_fd")
				if err != nil {
					logging.Fatal(&map[string]string{
						"file":     "authorization.go",
						"Function": "AuthorizationRuleMiddleware",
						"error":    "Can't parse userinfo from session",
					})
				}
				userinfo = userinfoI.(*util.FlatData)
			} else {
				logging.Fatal(&map[string]string{
					"file":     "authorization.go",
					"Function": "AuthorizationRuleMiddleware",
					"error":    "Can't get userinfo from session",
				})
			}
		}

		if !skip {
			// Check if request match required rule
			for _, rule := range requiredRules {
				math, allow := rule.TestRequired(r, userinfo, queryFD, bodyFD)

				if math {
					if !allow {
						logging.Debug(&map[string]string{
							"file":     "authorization.go",
							"Function": "AuthorizationRuleMiddleware",
							"event":    "Required rule violated - send forbidden",
							"rule":     fmt.Sprintf("%+v", rule),
						})

						http.Error(w, "403 - Forbidden.", 403)
						return
					}

					skip = true
					break
				}
			}
		}

		if !skip {
			// Check if request match disallow rule
			for _, rule := range disallowRules {
				math, allow := rule.TestDisallow(r, userinfo, queryFD, bodyFD)

				if math {
					if !allow {
						logging.Debug(&map[string]string{
							"file":     "authorization.go",
							"Function": "AuthorizationRuleMiddleware",
							"event":    "Disallow rule violated - send forbidden",
							"rule":     fmt.Sprintf("%+v", rule),
						})

						http.Error(w, "403 - Forbidden.", 403)
						return
					}

					skip = true
					break
				}
			}
		}

		// Serve to upstream
		logging.Debug(&map[string]string{
			"file":     "authorization.go",
			"Function": "AuthorizationRuleMiddleware",
			"event":    "Serve to upstream",
		})
		next.ServeHTTP(w, r)
	})
}
