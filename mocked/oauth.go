package mocked

import (
	"auth-guardian/config"
	"auth-guardian/util"
	"encoding/json"
	"fmt"
	"net/http"
)

var scheme string
var clientID string
var clientSecret string
var authCode string

// RunMockOAuthIDP runs a mocked OAuth IDP
func RunMockOAuthIDP() {
	scheme = "http"
	if config.IsHTTPS {
		scheme = "https"
	}

	clientID = "See you space"
	clientSecret = "cowboy"
	authCode = "SomeRandomStuff8heDDefDAwt83"

	mux := http.NewServeMux()

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		queryMap := map[string]interface{}{}
		for key, value := range r.URL.Query() {
			queryMap[key] = value
		}
		queryFD := util.NewFlatData()
		queryFD.BuildFrom(queryMap)

		// Check client id
		_clientID, err := queryFD.Search("client_id")
		if err != nil || _clientID.([]string)[0] != clientID {
			http.Error(w, "Invalid client id", http.StatusBadRequest)
			return
		}

		// Check request type
		responseType, err := queryFD.Search("response_type")
		if err != nil || responseType.([]string)[0] != "code" {
			http.Error(w, "Responsetype not code", http.StatusBadRequest)
			return
		}

		// Get state from request
		_state, err := queryFD.Search("state")
		if err != nil {
			http.Error(w, "State not sended", http.StatusBadRequest)
			return
		}
		state := _state.([]string)[0]

		// Redirect user to callback
		http.Redirect(w, r, fmt.Sprintf("%s://localhost:3000/callback?code=%s&state=%s", scheme, authCode, state), http.StatusFound)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// Send access token
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=mocktoken&scope=user,email&token_type=bearer"))
	})

	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		// Check for authorization token
		barer := r.Header.Get("Authorization")
		if barer != "Bearer mocktoken" {
			http.Error(w, fmt.Sprintf("Invalid barere - Get \"%s\" want \"Bearer mocktoken\"", barer), http.StatusBadRequest)
			return
		}

		// Create fake userinfo
		userinfo := map[string]interface{}{
			"username": "Saitama",
			"email":    "saitama@punch.co.jp",
			"role":     []string{"Superhero", "Bald-Head", "Main-Protagonist"},
		}

		// Send userinfo back
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userinfo)
	})

	http.ListenAndServe(":3002", mux)
}
