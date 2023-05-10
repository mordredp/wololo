package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Identify retrieves a session or creates a new one
func Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// We can obtain the session token from the requests cookies
		c, err := r.Cookie("session_token")

		var sessionToken string

		switch err {
		case nil:
			sessionToken = c.Value

		case http.ErrNoCookie:
			// if the cookie is not set
			// we create a new random session token
			// we use the "github.com/google/uuid" library to generate UUIDs
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(maxSessionLength)

			// we set the client cookie for "session_token"
			// as the session token we just generated
			// we also set an expiry time
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := userOrNew(sessionToken)

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
