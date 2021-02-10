package main

import (
	"auth-guardian/authmiddleware"
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/mocked"
	"auth-guardian/rules"
	"auth-guardian/upstream"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Load config
	err := config.Load()
	if err != nil {
		log.Panic(&map[string]string{
			"file":     "file.go",
			"Function": "getConfigFromFile",
			"error":    "Can't load existing config file",
			"details":  err.Error(),
		})
	}

	// Run test service if test mode enabled
	if config.MokeTestService {
		// Overrite config
		config.Upstream = "http://localhost:3001"

		// Run test service
		go mocked.Run()
	}

	// Run mock IDP if enabled
	if config.MokeOAuth {
		config.ClientID = "See you space"
		config.ClientSecret = "cowboy"
		config.AuthURL = "http://localhost:3002/auth"
		config.TokenURL = "http://localhost:3002/token"
		config.UserinfoURL = "http://localhost:3002/userinfo"

		go mocked.RunMockOAuthIDP()
	}

	// Initialize rules middleware
	rules.InitializeWhitelistMiddleware()

	// Initialize authorization middleware
	rules.InitializeAuthorizationMiddleware()

	// Initialize the OAuth middleware
	authInit, authMiddleware := authmiddleware.Provide()
	authInit()

	// Initialize upstream
	logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Initialize reverse proxy"})
	upstreamInit, upstreamProxy := upstream.Provide()
	// Initialize reverse proxy
	upstreamInit()
	// Set request handler
	http.Handle("/",
		rules.WhitelistRuleMiddleware(
			authMiddleware(rules.AuthorizationRuleMiddleware(upstreamProxy())),
			upstreamProxy(),
		))

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
