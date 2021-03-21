package server

import (
	"auth-guardian/authmiddleware"
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/rules"
	"auth-guardian/upstream"
	"fmt"
	"log"
	"net/http"
	"testing"
)

// Run the server
func Run(t *testing.T) *http.Server {
	// Initialize rules middleware
	rules.InitializeWhitelistMiddleware()

	// Initialize authorization middleware
	rules.InitializeAuthorizationMiddleware()

	// Initialize the middleware
	authInit, authMiddleware := authmiddleware.Provide()
	authInit()

	// Initialize upstream
	logging.Debug(&map[string]string{"file": "main.go", "Function": "main", "event": "Initialize reverse proxy"})
	upstreamInit, upstreamProxy := upstream.Provide()

	// Initialize reverse proxy
	upstreamInit()

	// Set request handler
	mux := http.NewServeMux()
	mux.Handle("/",
		rules.WhitelistRuleMiddleware(
			authMiddleware(rules.AuthorizationRuleMiddleware(upstreamProxy())),
			upstreamProxy(),
		),
	)

	// Run server in test mode
	if t != nil {
		server := &http.Server{Addr: config.Listen, Handler: mux}

		go func() {
			if config.IsHTTPS {
				// Run HTTPS server
				if err := server.ListenAndServeTLS(config.ServerCrt, config.ServerKey); err != http.ErrServerClosed {
					t.Errorf("Server statup failed, reason: %s", err.Error())
				}
			} else {
				// Run HTTP server
				if err := server.ListenAndServe(); err != http.ErrServerClosed {
					t.Errorf("Server statup failed, reason: %s", err.Error())
				}
			}
		}()

		return server
	}

	// Run server
	logging.Info(&map[string]string{"event": fmt.Sprintf("Listening on %s...", config.Listen)})

	if config.IsHTTPS {
		// Run HTTPS server
		if err := http.ListenAndServeTLS(config.Listen, config.ServerCrt, config.ServerKey, mux); err != nil {
			log.Fatalf("Server statup failed, reason: %s", err.Error())
		}
	} else {
		// Run HTTP server
		if err := http.ListenAndServe(config.Listen, mux); err != nil {
			log.Fatalf("Server statup failed, reason: %s", err.Error())
		}
	}

	return nil
}
