package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Login authenticates the session assigned to a user
func Login(w http.ResponseWriter, r *http.Request) {

	creds := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	expectedPassword := "code"

	if creds.Password != expectedPassword {

		//w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(maxSessionLength)

	sessions[sessionToken] = session{
		ID:     creds.Username,
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	log.Printf("user \"%s\" logged in", creds.Username)

	http.Redirect(w, r, "/", http.StatusFound)
}
