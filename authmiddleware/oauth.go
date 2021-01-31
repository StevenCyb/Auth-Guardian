package authmiddleware

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var oauthConfig *oauth2.Config

// InitOAuthMiddleware initialize the OAuth middleware
func InitOAuthMiddleware() {
	logging.Debug(&map[string]string{"file": "oauth.go", "Function": "InitOAuthMiddleware", "event": "Initialize OAuth middleware"})

	oauthConfig = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}

	go OAuthSessionRefresh()
}

// OAuthSessionRefresh refresh token and userinfo
func OAuthSessionRefresh() {
	// Lock session map
	SessionMapLock.Lock()

	logging.Debug(&map[string]string{
		"file":     "oauth.go",
		"Function": "OAuthSessionRefresh",
		"event":    "Run OAuth session refresh job",
	})

	for _, session := range SessionMap {
		if sat, err := session.Get("access_token"); err == nil {
			token := sat.(*oauth2.Token)
			if token.RefreshToken != "" && token.Expiry.Add(-1*time.Minute).Unix() <= time.Now().Unix() {
				// Refresh the token
				refreshToken(token)

				// Store new access token in session
				session.Set("access_token", token)

				// Get claims from access token and set it to session if not error
				parseClaims(session, token)

				// Get userinfo if URL defined
				if config.UserinfoURL != "" {
					getUserinfo(session, token)
				}
			}
		}
	}

	// Unlock session map
	SessionMapLock.Unlock()

	// Wait one minute to repeat job
	time.AfterFunc(1*time.Minute, func() { OAuthSessionRefresh() })
}

// OAuthMiddleware is the middleware to provide OAuth mechanism
func OAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug(&map[string]string{
			"file":       "oauth.go",
			"Function":   "OAuthMiddleware",
			"event":      "Handle request",
			"method":     r.Method,
			"req_scheme": r.URL.Scheme,
			"req_host":   r.Host,
			"req_path":   r.URL.Path,
			"req_query":  r.URL.RawQuery,
		})

		// Start a session
		session := SessionStart(w, r)

		// Check if session has access token, if yes serve to next handler
		if sat, err := session.Get("access_token"); err == nil && fmt.Sprint(sat.(*oauth2.Token).AccessToken) != "" {
			logging.Debug(&map[string]string{"file": "oauth.go", "Function": "OAuthMiddleware", "event": "Forward"})
			next.ServeHTTP(w, r)
			return
		}

		// If its a callback
		if r.URL.Path == "/callback" {
			// Get state
			state, err := GetCookieValue(r, "OG_STATE")
			if err != nil {
				logging.Error(&map[string]string{
					"file":     "oauth.go",
					"Function": "OAuthMiddleware",
					"name":     "OG_STATE",
					"error":    err.Error(),
				})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Destroy state cookie
			DestroyCookie(w, "OG_STATE", http.SameSiteLaxMode)
			logging.Debug(&map[string]string{
				"file":     "oauth.go",
				"Function": "OAuthMiddleware",
				"event":    "Handle callback",
				"state":    state,
			})

			// Check if state is
			if r.FormValue("state") != state {
				logging.Error(&map[string]string{
					"file":     "oauth.go",
					"Function": "OAuthMiddleware",
					"error":    "Invalid OAuth state got:" + r.FormValue("state") + " want:" + state,
				})
				if w != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			// Request access token
			token, err := oauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
			if err != nil {
				logging.Error(&map[string]string{
					"file":     "oauth.go",
					"Function": "OAuthMiddleware",
					"error":    "Code exchange wrong:" + err.Error(),
				})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !token.Valid() {
				logging.Error(&map[string]string{
					"file":     "oauth.go",
					"Function": "OAuthMiddleware",
					"error":    "Retrieved token is invalid ",
				})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Store access token in session
			session.Set("access_token", token)

			// Get claims from access token and set it to session if not error
			parseClaims(session, token)

			// Get userinfo if URL defined
			if config.UserinfoURL != "" {
				getUserinfo(session, token)
			}

			// Construct a redirect URl based on the original request
			redirectPath := r.Host
			if config.IsHTTPS {
				redirectPath = "https://" + redirectPath
			} else {
				redirectPath = "http://" + redirectPath
			}
			if value, err := session.Get("path_before_redirect"); err == nil {
				redirectPath = redirectPath + value.(string)
			} else {
				redirectPath = redirectPath + "/"
			}

			// Redirect from client side to keep the cookies (including session id)
			logging.Debug(&map[string]string{
				"file":         "oauth.go",
				"Function":     "OAuthMiddleware",
				"event":        "Client side redirect",
				"redirect_url": redirectPath,
			})
			fmt.Fprintf(w, "<html><head></head><body onload=\"window.location = '%s';\"></body><html>", redirectPath)
			return
		}

		// Store requested url
		session.Set("path_before_redirect", fmt.Sprintf("%v", r.URL))

		// Redirect to auth url
		state := GetRandomBase64String(16)
		SetCookie(w, "OG_STATE", state, http.SameSiteLaxMode, time.Duration(config.StateLifetime))

		url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
		logging.Debug(&map[string]string{
			"file":     "oauth.go",
			"Function": "OAuthMiddleware",
			"event":    "Redirect to auth",
			"url":      url,
		})

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// Parse claims of token and set to session
func parseClaims(session *Session, token *oauth2.Token) {
	var claims map[string]interface{}
	tokenClaimString, err := jwt.ParseSigned(token.AccessToken)
	if err == nil {
		_ = tokenClaimString.UnsafeClaimsWithoutVerification(&claims)
		session.Set("claims", claims)
	}
}

// Get userinfo and set to session
func getUserinfo(session *Session, token *oauth2.Token) {
	logging.Debug(&map[string]string{
		"file":     "oauth.go",
		"Function": "OAuthMiddleware",
		"event":    "Get userinfo",
		"url":      config.UserinfoURL,
	})

	// Request userinfo
	client := oauthConfig.Client(oauth2.NoContext, token)
	response, err := client.Get(config.UserinfoURL)
	if err != nil {
		logging.Warning(&map[string]string{
			"file":     "oauth.go",
			"Function": "OAuthMiddleware",
			"warning":  "Failed getting user info:" + err.Error(),
		})
	}
	defer response.Body.Close()

	// Read userinfo
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Warning(&map[string]string{
			"file":     "oauth.go",
			"Function": "OAuthMiddleware",
			"warning":  "Failed read response:" + err.Error(),
		})
	}

	// Parse userinfo
	userinfo, err := JSONToMap(string(body))
	if err != nil {
		logging.Warning(&map[string]string{
			"file":     "oauth.go",
			"Function": "OAuthMiddleware",
			"warning":  "Failed to map userinfo:" + err.Error(),
		})
	}
	session.Set("userinfo", userinfo)
}

// Request a new token
func refreshToken(token *oauth2.Token) {
	logging.Debug(&map[string]string{
		"file":     "oauth.go",
		"Function": "refreshToken",
		"event":    "Refresh token in session",
	})

	// Check if refresh token exists
	if token.RefreshToken == "" {
		return
	}

	// Build request of new access token
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("client_id", oauthConfig.ClientID)
	v.Set("client_secret", oauthConfig.ClientSecret)
	v.Set("refresh_token", token.RefreshToken)

	req, err := http.NewRequest("POST", config.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Do request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		fmt.Printf("cannot fetch token: %v\n", err.Error())
		return
	}

	// Parse response body
	content, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if content == "application/x-www-form-urlencoded" || content == "text/plain" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		token = &oauth2.Token{
			AccessToken:  values.Get("access_token"),
			TokenType:    values.Get("token_type"),
			RefreshToken: values.Get("refresh_token"),
		}

		e := values.Get("expires_in")
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	} else {
		tj := struct {
			AccessToken  string
			TokenType    string
			RefreshToken string
			ExpiresIn    int32
		}{}

		if err = json.Unmarshal(body, &tj); err != nil {
			fmt.Println(err.Error())
			return
		}

		token = &oauth2.Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       time.Now().Add(time.Duration(tj.ExpiresIn) * time.Second),
		}
	}
}
