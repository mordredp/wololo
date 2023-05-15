package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Login authenticates the session assigned to a user.
func (a *Authenticator) Login(w http.ResponseWriter, r *http.Request) {

	c := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if err := a.provider.Authenticate(c.Username, c.Password); err != nil {
		log.Printf("provider error: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(a.maxSessionLength)

	a.sessions[sessionToken] = session{
		ID:     c.Username,
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	log.Printf("user %q logged in", c.Username)

	http.Redirect(w, r, "/", http.StatusFound)
}
