package rules

import (
	"auth-guardian/logging"
	"net/http"
)

// AuthorizationRuleMiddleware check request are authorized or unauthorized through rules
func AuthorizationRuleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "authorization.go",
			"Function":   "AuthorizationRuleMiddleware",
			"event":      "Handle request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.Path,
			"req_query":  r.URL.RawQuery,
		})

		// Serve
		next.ServeHTTP(w, r)
	})
}
