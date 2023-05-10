package auth

import (
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

	// create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(maxSessionLength)

	// set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		ID:     creds.Username,
		expiry: expiresAt,
	}

	// we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
