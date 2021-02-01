package main

import (
	"auth-guardian/authmiddleware"
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/testservice"
	"auth-guardian/upstream"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Load config
	version := config.Load()
	if version {
		fmt.Println("Version: 0.1.1")
		return
	}

	// Run test service if test mode enabled
	if config.TestMode {
		// Overrite config
		config.Upstream = "http://localhost:3001"

		// Run test service
		go testservice.Run()
	}

	// Initialize the OAuth middleware
	authmiddleware.InitOAuthMiddleware()

	// Initialize upstream
	if !config.CORSUpstream {
		// Initialize default reverse proxy
		logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Use default ReverseProxy"})
		upstream.InitDefaultReverseProxy()

		// Set request handler
		logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Set default reverse proxy handler"})
		http.Handle("/", authmiddleware.OAuthMiddleware(upstream.ProxyHandler()))
	} else {
		// Initialize CORS reverse proxy
		logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Use CORS ReverseProxy"})
		upstream.InitCORSReverseProxy()

		// Set request handler
		logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Set CORS reverse proxy handler"})
		http.Handle("/", authmiddleware.OAuthMiddleware(upstream.CORSProxyHandler()))
	}

	// Run server
	logging.Info(&map[string]string{"event": fmt.Sprintf("Listening on %s...", config.Listen)})
	if config.IsHTTPS {
		// Run HTTPS server
		err := http.ListenAndServeTLS(config.Listen, config.ServerCrt, config.ServerKey, nil)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// Run HTTP server
		err := http.ListenAndServe(config.Listen, nil)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
