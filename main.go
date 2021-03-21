package main

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/mocked"
	"auth-guardian/server"
	"context"
	"log"
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
	if config.MockTestService {
		logging.Debug(&map[string]string{
			"file":     "main.go",
			"Function": "main",
			"event":    "Run mocked test-service",
			"url":      "http://localhost:3001",
		})

		// Run test service
		testServiceServer := mocked.RunMockTestService()
		defer testServiceServer.Shutdown(context.TODO())
	}

	// Run mock IDP if enabled
	if config.MockOAuth {
		logging.Debug(&map[string]string{
			"file":     "main.go",
			"Function": "main",
			"event":    "Run mocked OAuth IDP",
		})

		// Run OAuth IDP
		oAuthServer := mocked.RunMockOAuthIDP()
		defer oAuthServer.Shutdown(context.TODO())
	} else if config.MockLDAP {
		logging.Debug(&map[string]string{
			"file":     "main.go",
			"Function": "main",
			"event":    "Run mocked LDAP IDP",
		})

		// Run LDAP IDP
		ldapListener := mocked.RunMockLDAPIDP()
		defer (*ldapListener).Close()
	} else if config.MockSAML {
		logging.Debug(&map[string]string{
			"file":     "main.go",
			"Function": "main",
			"event":    "Run mocked SAML IDP",
		})

		// Run SAML IDP
		samlServer := mocked.RunMockSAMLIDP()
		defer samlServer.Shutdown(context.TODO())
	}

	server.Run(nil)
}
