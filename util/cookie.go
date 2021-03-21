package util

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"fmt"
	"net/http"
	"time"
)

var seed string

func init() {
	seed = "_" + GetRandomString(8)
}

// SetCookie create and set a cookie with given name, value and lifetime (min)
func SetCookie(w http.ResponseWriter, name string, value string, sameSite http.SameSite, minuteLifetime time.Duration) {
	logging.Debug(&map[string]string{
		"file":      "cookie.go",
		"Function":  "SetCookie",
		"event":     "Create cookie",
		"name":      name,
		"value":     value,
		"same_site": fmt.Sprint(sameSite),
		"lifetime":  fmt.Sprint(minuteLifetime * time.Minute),
		"seed":      seed,
	})

	cookie := http.Cookie{
		Name:     name + seed,
		Value:    value,
		Expires:  time.Now().Add(minuteLifetime * time.Minute),
		Path:     "/",
		SameSite: sameSite,
		Secure:   config.IsHTTPS,
	}
	http.SetCookie(w, &cookie)
}

// DestroyCookie destroys a cookie with name
func DestroyCookie(w http.ResponseWriter, name string, sameSite http.SameSite) {
	logging.Debug(&map[string]string{
		"file":      "cookie.go",
		"Function":  "DestroyCookie",
		"event":     "Destroy cookie",
		"name":      name,
		"same_site": fmt.Sprint(sameSite),
		"seed":      seed,
	})
	newCookie := &http.Cookie{
		Name:     name + seed,
		Expires:  time.Now().Add(-1 * time.Second),
		Path:     "/",
		SameSite: sameSite,
		Secure:   config.IsHTTPS,
	}

	http.SetCookie(w, newCookie)
}

// GetCookieValue return value of a cookie with name
func GetCookieValue(r *http.Request, name string) (string, error) {
	logging.Debug(&map[string]string{
		"file":     "cookie.go",
		"Function": "GetCookieValue",
		"event":    "Get cookie",
		"name":     name,
		"seed":     seed,
	})

	cookie, err := r.Cookie(name + seed)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
