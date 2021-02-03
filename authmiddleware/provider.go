package authmiddleware

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"net/http"
)

// Provide give the configured upstream reference
func Provide() (func(), func(next http.Handler) http.Handler) {
	if validOAuthConfiguration() {
		logging.Debug(&map[string]string{"file": "authmiddleware/provider.go", "Function": "Provide", "event": "Provide OAuth authentication"})
		return InitOAuthMiddleware, OAuthMiddleware
	} else if validSAMLConfiguration() {
		logging.Debug(&map[string]string{"file": "authmiddleware/provider.go", "Function": "Provide", "event": "Provide SAML authentication"})
		return InitSAMLhMiddleware, SAMLMiddleware
	}

	logging.Fatal(&map[string]string{"file": "authmiddleware/provider.go", "Function": "Provide", "error": "Configuration not match any of provided authentication mechanisms"})
	return nil, nil
}

func validOAuthConfiguration() bool {
	return config.ClientID != "" && config.ClientSecret != "" && config.AuthURL != "" && config.TokenURL != ""
}

func validSAMLConfiguration() bool {
	return config.SAMLKey != "" && config.SAMLCrt != "" && config.IdpMetadataURL != "" && config.SelfRootURL != ""
}
