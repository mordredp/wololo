package auth

import (
	"log"
	"net/http"
	"time"
)

// Logout removes a session.
func (a *Authenticator) Logout(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			//w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		//w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
	sessionToken := c.Value

	delete(a.sessions, sessionToken)

	log.Printf("user %q logged out", r.Context().Value(UserKey).(User).ID)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
