package upstream

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"auth-guardian/util"
	"fmt"
	"net/http"
)

// Set configured forward information's in cookie
func setForwardInformations(sourceReq *http.Request, targetReq *http.Request) {
	if config.ForwardAccessToken || config.ForwardUserinfo {
		if sessionID := sourceReq.Context().Value(util.ContextKey("session_id")); sessionID != nil {
			logging.Debug(&map[string]string{
				"file":         "upstream/utils.go",
				"Function":     "setForwardInformations",
				"event":        "Populate forwarded request",
				"access_token": fmt.Sprint(config.ForwardAccessToken),
				"userinfo":     fmt.Sprint(config.ForwardUserinfo),
			})

			session := util.SessionMap[sessionID.(string)]
			if config.ForwardUserinfo {
				userinfo, err := session.Get("userinfo_string")
				if err == nil {
					targetReq.AddCookie(&http.Cookie{Name: "userinfo", Value: userinfo.(string)})
				} else {
					logging.Warning(&map[string]string{
						"file":     "upstream/utils.go",
						"Function": "setForwardInformations",
						"warning":  "Userinfo not in session",
					})
				}
			}
			if config.ForwardAccessToken {
				accessToken, err := session.Get("access_token_string")
				if err == nil {
					targetReq.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken.(string)})
				} else {
					logging.Warning(&map[string]string{
						"file":     "upstream/utils.go",
						"Function": "setForwardInformations",
						"warning":  "Access token not in session",
					})
				}
			}
		} else {
			logging.Warning(&map[string]string{
				"file":     "upstream/utils.go",
				"Function": "setForwardInformations",
				"warning":  "Session ID not in context",
			})
		}
	}
}
