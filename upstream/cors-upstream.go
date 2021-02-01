package upstream

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

// InitCORSReverseProxy create a CORS reverse proxy
func InitCORSReverseProxy() {
	logging.Debug(&map[string]string{"file": "cors-upstream.go", "Function": "InitCORSReverseProxy", "event": "Initialize upstream"})

	Origin, _ = url.Parse(config.Upstream)
}

// CORSProxyHandler return a handler for CORS proxy
func CORSProxyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set forward header
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", Origin.Host)

		// Copy request and modify it by setting the upstream
		nextReq := r.Clone(context.TODO())
		nextReq.RequestURI = ""
		nextReq.Host = Origin.Host
		nextReq.URL.Host = Origin.Host
		nextReq.URL.Scheme = Origin.Scheme

		// Set referer
		if nextReq.Header.Get("referer") != "" {
			nextReq.Header.Set("referer", r.Header.Get("referer"))
		}

		logging.Debug(&map[string]string{
			"file":              "cors-upstream.go",
			"Function":          "CORSProxyHandler",
			"event":             "Forward request",
			"req_method":        r.Method,
			"req_scheme":        r.URL.Scheme,
			"req_host":          r.Host,
			"req_path":          r.URL.RawPath,
			"req_query":         r.URL.RawQuery,
			"forward_method":    nextReq.Method,
			"forward_scheme":    nextReq.URL.Scheme,
			"forward_req_host":  nextReq.Host,
			"forward_req_path":  nextReq.URL.RawPath,
			"forward_req_query": nextReq.URL.RawQuery,
		})

		// Set configured forward information's in cookie
		setForwardInformations(r, nextReq)

		// Do request
		client := http.Client{}
		nextRes, err := client.Do(nextReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer nextRes.Body.Close()

		// Set header
		w.Header().Set("Access-Control-Allow-Origin", "*")
		for key, valueSlice := range nextRes.Header {
			if key == "Access-Control-Allow-Origin" {
				continue
			}
			for _, value := range valueSlice {
				w.Header().Add(key, value)
			}
		}

		// Set body content and send response back to requester
		body, err := ioutil.ReadAll(nextRes.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
