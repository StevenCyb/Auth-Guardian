package authmiddleware

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/util"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/ldap"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator

// InitLDAPhMiddleware initialize the LDAP middleware
func InitLDAPhMiddleware() {
	logging.Debug(&map[string]string{"file": "ldap.go", "Function": "InitSAMLhMiddleware", "event": "Initialize SAML middleware"})

	// Create configuration
	cfg := &ldap.Config{
		BaseDN:       config.DirectoryServerBaseDN,
		BindDN:       config.DirectoryServerBindDN,
		Port:         fmt.Sprint(config.DirectoryServerPort),
		Host:         config.DirectoryServerHost,
		BindPassword: config.DirectoryServerBindPassword,
		Filter:       config.DirectoryServerFilter,
	}

	// Create and init LDAP authenticator
	authenticator = auth.New()
	cache := store.NewFIFO(context.Background(), time.Minute*time.Duration(config.SessionLifetime))
	strategy := ldap.NewCached(cfg, cache)
	authenticator.EnableStrategy(ldap.StrategyKey, strategy)
}

// LDAPMiddleware is the middleware to provide LDAP authentication mechanism
func LDAPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "ldap.go",
			"Function":   "LDAPMiddleware",
			"event":      "Handle request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.Path,
			"req_query":  r.URL.RawQuery,
		})

		// Start own session implementation
		session := util.SessionStart(w, r)

		// Check if userinfo in session
		_, userinfoError := session.Get("userinfo_string")

		// Check if user authenticated
		user, err := authenticator.Authenticate(r)
		if err != nil {
			logging.Debug(&map[string]string{"file": "ldap.go", "Function": "LDAPMiddleware", "event": "Send unauthorized"})

			// Clear session if has userinfo
			if userinfoError == nil {
				logging.Debug(&map[string]string{"file": "ldap.go", "Function": "LDAPMiddleware", "event": "Clear session before new authentication"})
				util.SessionDestroy(session.SID)
			}

			// send unauthorized back
			w.Header().Set("WWW-Authenticate", `Basic realm="Auth-Guardian"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		// Set userinfo if not present
		if userinfoError != nil {
			logging.Debug(&map[string]string{"file": "ldap.go", "Function": "LDAPMiddleware", "event": "Set userinfo to new session"})

			// Store userinfo as base64 string to forward if configured
			userinfoString, _ := json.Marshal(user)
			session.Set("userinfo_string", base64.StdEncoding.EncodeToString(userinfoString))

			userinfoFlatData := util.NewFlatData()
			userinfoFlatData.BuildFrom(map[string]interface{}{
				"username":   user.UserName(),
				"user_id":    user.ID(),
				"groups":     user.Groups(),
				"extensions": user.Extensions(),
			})
			session.Set("userinfo_fd", userinfoFlatData)
		}

		// Set session id to context
		ctx := context.WithValue(r.Context(), util.ContextKey("session_id"), session.SID)

		// Serve to next handler with context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
