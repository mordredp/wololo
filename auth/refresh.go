package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Refresh renews the session token expiration for valid tokens.
func (a *Authenticator) Refresh(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	userSession, exists := a.sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(a.sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(a.maxSessionLength)

	a.sessions[newSessionToken] = session{
		ID:     userSession.ID,
		expiry: expiresAt,
	}

	delete(a.sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(a.maxSessionLength),
	})
}
