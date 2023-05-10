package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Refresh renews the session token expiration for valid tokens
func Refresh(w http.ResponseWriter, r *http.Request) {

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

	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if the previous session is valid,
	// create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(maxSessionLength)

	// set the token in the session map,
	// along with the user whom it is assigned to
	sessions[newSessionToken] = session{
		ID:     userSession.ID,
		expiry: expiresAt,
	}

	// delete the older session token
	delete(sessions, sessionToken)

	// set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(maxSessionLength),
	})
}
