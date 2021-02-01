package util

import (
	"auth-guardian/config"
	"auth-guardian/logging"
	"errors"
	"net/http"
	"sync"
	"time"
)

// SessionMapLock to lock the session map
var SessionMapLock sync.Mutex

// SessionMap store the sessions
var SessionMap map[string]*Session

// SessionGC is a garbage collector to remove old sessions
func SessionGC() {
	// Lock and defer unlock of session map
	SessionMapLock.Lock()
	defer SessionMapLock.Unlock()

	logging.Debug(&map[string]string{
		"file":     "session.go",
		"Function": "SessionGC",
		"event":    "Run garbage collector job",
	})

	// Destory session if last access older than configured session lifetime
	for SID, value := range SessionMap {
		if value.TimeAccessed.Add(time.Duration(config.SessionLifetime)*time.Minute).Unix() <= time.Now().Unix() {
			SessionDestroy(SID)
		}
	}

	// Wait one minute to repeat job
	time.AfterFunc(time.Duration(config.SessionLifetime+1)*time.Minute, func() { SessionGC() })
}

func init() {
	SessionMap = make(map[string]*Session)
	go SessionGC()
}

// SessionExists check if session exist on serverside
func SessionExists(SID string) bool {
	_, exists := SessionMap[SID]
	return exists
}

// SessionStart start a new or use an existing session
func SessionStart(w http.ResponseWriter, r *http.Request) *Session {
	SessionMapLock.Lock()
	defer SessionMapLock.Unlock()

	SID, err := GetCookieValue(r, "OG_SESSION_ID")
	if err != nil || SID == "" || !SessionExists(SID) {
		oldSID := SID
		taken := true
		for taken {
			SID = GetRandomString(64)
			_, taken = SessionMap[SID]
		}

		logging.Debug(&map[string]string{
			"file":     "session.go",
			"Function": "SessionStart",
			"event":    "Create session",
			"old_SID":  oldSID,
			"new_SID":  SID,
		})

		SetCookie(w, "OG_SESSION_ID", SID, http.SameSiteStrictMode, time.Duration(config.SessionLifetime))

		session := Session{
			SID:          SID,
			TimeAccessed: time.Now(),
			Values:       make(map[string]interface{}),
		}

		SessionMap[SID] = &session
		return &session
	}

	logging.Debug(&map[string]string{
		"file":     "session.go",
		"Function": "SessionStart",
		"event":    "Use existing session",
		"SID":      SID,
	})
	SessionMap[SID].Used(w)
	return SessionMap[SID]
}

// SessionDestroy deletes a session
func SessionDestroy(SID string) {
	SessionMapLock.Lock()
	defer SessionMapLock.Unlock()

	logging.Debug(&map[string]string{
		"file":     "session.go",
		"Function": "SessionDestroy",
		"event":    "Destroy session",
		"SID":      SID,
	})
	delete(SessionMap, SID)
}

// Session object
type Session struct {
	SID          string
	TimeAccessed time.Time
	Values       map[string]interface{}
}

// Used set the last access time
func (s *Session) Used(w http.ResponseWriter) {
	SetCookie(w, "OG_SESSION_ID", s.SID, http.SameSiteStrictMode, time.Duration(config.SessionLifetime))
	s.TimeAccessed = time.Now()
}

// Set a new key value pair to session
func (s *Session) Set(key string, value interface{}) {
	s.Values[key] = value
}

// Get a value by key
func (s *Session) Get(key string) (interface{}, error) {
	if value, ok := s.Values[key]; ok {
		return value, nil
	}
	return "", errors.New("Item not exists in session")
}

// Delete a value by key
func (s *Session) Delete(key string) {
	delete(s.Values, key)
}
