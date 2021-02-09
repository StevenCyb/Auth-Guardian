package rules

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"fmt"
	"net/http"
)

var whitelistRules []Rule

// InitializeWhitelistMiddleware initialize the whitelist middleware
func InitializeWhitelistMiddleware() {
	logging.Debug(&map[string]string{"file": "whitelist.go", "Function": "InitializeWhitelistMiddleware", "event": "Initialize whitelist middleware"})

	for _, rule := range config.Rules {
		if rule.Type != "whitelist" {
			continue
		}
		ruleStruct := Rule{}
		ruleStruct.FromRuleConfig(rule)
		whitelistRules = append(whitelistRules, ruleStruct)

		logging.Info(&map[string]string{
			"event":       "Whitelist rule added",
			"rule_method": fmt.Sprintf("%v", rule.Method),
			"rule_path":   rule.Path,
		})
	}
}

// WhitelistRuleMiddleware check request again whitelist rules
func WhitelistRuleMiddleware(nextSecured http.Handler, nextWhitelisted http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":      "whitelist.go",
			"Function":  "WhitelistRuleMiddleware",
			"event":     "Handle request",
			"method":    r.Method,
			"req_path":  r.URL.Path,
			"req_query": r.URL.RawQuery,
		})

		// Check if request match whitelist rule
		for _, rule := range whitelistRules {
			if rule.TestWhitelist(r) {
				logging.Debug(&map[string]string{
					"file":     "whitelist.go",
					"Function": "WhitelistRuleMiddleware",
					"event":    "Serve to upstream",
				})
				nextWhitelisted.ServeHTTP(w, r)
				return
			}
		}

		// Serve to auth middleware
		logging.Debug(&map[string]string{
			"file":     "whitelist.go",
			"Function": "WhitelistRuleMiddleware",
			"event":    "Serve to auth middleware",
		})
		nextSecured.ServeHTTP(w, r)
	})
}
