package upstream

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy for upstream server
var ReverseProxy *httputil.ReverseProxy

// Origin of the upstream
var Origin *url.URL

// InitDefaultReverseProxy create a default reverse proxy
func InitDefaultReverseProxy() {
	logging.Debug(&map[string]string{"file": "default-upstream.go", "Function": "InitDefaultReverseProxy", "event": "Initialize upstream"})

	Origin, _ = url.Parse(config.Upstream)
	ReverseProxy = httputil.NewSingleHostReverseProxy(Origin)
	ReverseProxy.Director = DirectorHandler
}

// DirectorHandler set origin and header
func DirectorHandler(r *http.Request) {
	logging.Debug(&map[string]string{
		"file":       "default-upstream.go",
		"Function":   "DirectorHandler",
		"event":      "Set origin and header",
		"req_host":   r.Host,
		"ori_host":   Origin.Host,
		"req_scheme": r.URL.Scheme,
		"ori_scheme": Origin.Scheme,
	})

	r.Host = Origin.Host
	r.URL.Host = Origin.Host
	r.URL.Scheme = Origin.Scheme

	// Set forward header
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Origin-Host", Origin.Host)
}

// DefaultProxyHandler return a handler for default proxy
func DefaultProxyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "default-upstream.go",
			"Function":   "ProxyHandler",
			"event":      "Forward request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.RawPath,
			"req_query":  r.URL.RawQuery,
		})

		// Set CORS header
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Set configured forward infromations in cookie
		setForwardInformations(r, r)

		ReverseProxy.ServeHTTP(w, r)
	})
}
