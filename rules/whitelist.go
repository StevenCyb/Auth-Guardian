package rules

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"fmt"
	"net/http"
	"regexp"
)

var whitelistRules []*regexp.Regexp

// InitializeWhitelistMiddleware initialize the whitelist middleware
func InitializeWhitelistMiddleware() {
	logging.Debug(&map[string]string{"file": "whitelist.go", "Function": "InitializeWhitelistMiddleware", "event": "Initialize whitelist middleware"})

	for _, rule := range config.WhitelistRules {
		reg, err := regexp.Compile(rule)
		if err != nil {
			logging.Fatal(&map[string]string{
				"file":         "whitelist.go",
				"Function":     "InitializeWhitelistMiddleware",
				"event":        "Whitelist regex parsing failed",
				"regex_string": rule,
			})
		}
		whitelistRules = append(whitelistRules, reg)
	}
}

// WhitelistRuleMiddleware check request again whitelist rules
func WhitelistRuleMiddleware(nextSecured http.Handler, nextWhitelisted http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "whitelist.go",
			"Function":   "WhitelistRuleMiddleware",
			"event":      "Handle request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.Path,
			"req_query":  r.URL.RawQuery,
		})

		// Build request string
		requestString := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		// Check if request match whitelist rule
		for _, rule := range whitelistRules {
			if rule.MatchString(requestString) {
				nextWhitelisted.ServeHTTP(w, r)
				return
			}
		}

		// Serve to auth middleware
		nextSecured.ServeHTTP(w, r)
	})
}
