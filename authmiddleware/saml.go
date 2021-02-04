package authmiddleware

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/util"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crewjam/saml/samlsp"
)

var samlSP *samlsp.Middleware

// InitSAMLhMiddleware initialize the SAML middleware
func InitSAMLhMiddleware() {
	logging.Debug(&map[string]string{"file": "oauth.go", "Function": "InitSAMLhMiddleware", "event": "Initialize SAML middleware"})

	// Load key and certificate of this SP
	keyPair, err := tls.LoadX509KeyPair(config.SAMLCrt, config.SAMLKey)
	if err != nil {
		fmt.Println(err)
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "Loading of certificate failed"})
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "Parsing of certificate failed"})
	}

	// Fetch IDP metadata
	idpMetadataURL, err := url.Parse(config.IdpMetadataURL)
	if err != nil {
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "IDP metadata URL parsing failed"})
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "IDP metadata fetch failed"})
	}

	// Parse the URL to self
	rootURL, err := url.Parse(config.SelfRootURL)
	if err != nil {
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "Self root URL parsing failed"})
	}

	// Initialize SAML-SP
	samlSP, err = samlsp.New(samlsp.Options{
		URL:          *rootURL,
		Key:          keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:  keyPair.Leaf,
		IDPMetadata:  idpMetadata,
		SignRequest:  true,
		CookieMaxAge: 5 * time.Minute,
	})
	if err != nil {
		logging.Fatal(&map[string]string{"file": "authmiddleware/saml.go", "Function": "InitSAMLhMiddleware", "error": "Initialization of SAML-SP failed"})
	}
}

// SAMLMiddleware is the middleware to provide SAML mechanism
func SAMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "saml.go",
			"Function":   "SAMLMiddleware",
			"event":      "Handle request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.Path,
			"req_query":  r.URL.RawQuery,
		})

		// If its a callback
		if strings.HasPrefix(r.URL.Path, "/saml/") {
			samlSP.ServeHTTP(w, r)
			return
		}

		// Check if user has SAML session
		samlSession, err := samlSP.Session.GetSession(r)
		// Start own session implementation
		session := util.SessionStart(w, r)
		if err == nil && samlSession != nil && session != nil {
			// Forward to next if authenticated
			logging.Debug(&map[string]string{"file": "oauth.go", "Function": "OAuthMiddleware", "event": "Forward with context"})

			samlSessionWithAttributes, ok := samlSession.(samlsp.SessionWithAttributes)
			if ok {
				userinfo := samlSessionWithAttributes.GetAttributes()
				session.Set("userinfo", userinfo)

				userinfoString, _ := json.Marshal(userinfo)
				session.Set("userinfo_string", base64.StdEncoding.EncodeToString(userinfoString))
			} else {
				logging.Warning(&map[string]string{"file": "oauth.go", "Function": "OAuthMiddleware", "warning": "Can't get SAML session attributes"})
			}

			// Set session id to context
			ctx := context.WithValue(r.Context(), util.ContextKey("session_id"), session.SID)

			// Serve to next handler with context
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		samlSP.HandleStartAuthFlow(w, r)
	})
}
