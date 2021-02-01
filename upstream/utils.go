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
			if config.ForwardAccessToken {
				userinfo, err := session.Get("userinfo_string")
				if err == nil {
					targetReq.AddCookie(&http.Cookie{Name: "userinfo", Value: userinfo.(string)})
				} else {
					logging.Warning(&map[string]string{
						"file":     "upstream/utils.go",
						"Function": "setForwardInformations",
						"event":    "Userinfo not in session",
					})
				}
			}
			if config.ForwardUserinfo {
				accessToken, err := session.Get("access_token_string")
				if err == nil {
					targetReq.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken.(string)})
				} else {
					logging.Warning(&map[string]string{
						"file":     "upstream/utils.go",
						"Function": "setForwardInformations",
						"event":    "Access token not in session",
					})
				}
			}
		} else {
			logging.Warning(&map[string]string{
				"file":     "upstream/utils.go",
				"Function": "setForwardInformations",
				"event":    "Session ID not in context",
			})
		}
	}
}

// // Set configured forward information's in cookie
// func setForwardInformations(sourceReq *http.Request, targetReq *http.Request) {
// 	if config.ForwardAccessToken || config.ForwardUserinfo {
// 		if sessionID := sourceReq.Context().Value(util.ContextKey("session_id")); sessionID != nil {
// 			logging.Debug(&map[string]string{
// 				"file":         "upstream/utils.go",
// 				"Function":     "setForwardInformations",
// 				"event":        "Populate forwarded request",
// 				"access_token": fmt.Sprint(config.ForwardAccessToken),
// 				"userinfo":     fmt.Sprint(config.ForwardUserinfo),
// 			})

// 			session := util.SessionMap[sessionID.(string)]
// 			if config.ForwardAccessToken {
// 				userinfo, err := session.Get("userinfo")
// 				if err == nil {
// 					userinfoString, err := json.Marshal(userinfo)

// 					if err == nil && len(userinfoString) > 0 {
// 						targetReq.AddCookie(&http.Cookie{Name: "userinfo", Value: string(userinfoString)})
// 					} else {
// 						logging.Warning(&map[string]string{
// 							"file":     "upstream/utils.go",
// 							"Function": "setForwardInformations",
// 							"event":    "Cant marshal userinfo",
// 						})
// 					}
// 				} else {
// 					logging.Warning(&map[string]string{
// 						"file":     "upstream/utils.go",
// 						"Function": "setForwardInformations",
// 						"event":    "Userinfo not in session",
// 					})
// 				}
// 			}
// 			if config.ForwardUserinfo {
// 				accessToken, err := session.Get("access_token")
// 				if err == nil {
// 					accessTokenString, err := json.Marshal(accessToken)

// 					if err == nil && len(accessTokenString) > 0 {
// 						targetReq.AddCookie(&http.Cookie{Name: "access_token", Value: string(accessTokenString)})
// 					} else {
// 						logging.Warning(&map[string]string{
// 							"file":     "upstream/utils.go",
// 							"Function": "setForwardInformations",
// 							"event":    "Cant marshal access token",
// 						})
// 					}
// 				} else {
// 					logging.Warning(&map[string]string{
// 						"file":     "upstream/utils.go",
// 						"Function": "setForwardInformations",
// 						"event":    "Access token not in session",
// 					})
// 				}
// 			}
// 		} else {
// 			logging.Warning(&map[string]string{
// 				"file":     "upstream/utils.go",
// 				"Function": "setForwardInformations",
// 				"event":    "Session ID not in context",
// 			})
// 		}
// 	}
// }
