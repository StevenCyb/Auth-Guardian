package upstream

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"net/http"
)

// Provide give the configured upstream reference
func Provide() (func(), func() http.Handler) {
	if config.CORSUpstream {
		logging.Debug(&map[string]string{"file": "upstream/provider.go", "Function": "Provide", "event": "Use CORS ReverseProxy"})
		return InitCORSReverseProxy, CORSProxyHandler
	}

	logging.Debug(&map[string]string{"file": "upstream/provider.go", "Function": "Provide", "event": "Use default ReverseProxy"})
	return InitDefaultReverseProxy, DefaultProxyHandler
}
