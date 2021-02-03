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
	authInit, authMiddleware := authmiddleware.Provide()
	authInit()

	// Initialize upstream
	logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Initialize reverse proxy"})
	upstreamInit, upstreamProxy := upstream.Provide()
	// Initialize reverse proxy
	upstreamInit()
	// Set request handler
	http.Handle("/", authMiddleware(upstreamProxy()))

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
