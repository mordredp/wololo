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

		c, err := r.Cookie("session_token")

		var sessionToken string

		switch err {
		case nil:
			sessionToken = c.Value

		case http.ErrNoCookie:
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(maxSessionLength)

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
