package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type key int

const (
	// UserKey is the key to the value of a User in a context.
	UserKey key = iota
)

// User holds a users account information and its authentication status.
type User struct {
	ID            string
	Authenticated bool
}

// Credentials holds the username and password used to authenticate a session.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// a session contains an identifier (usually the username of the user
// it is assigned to) and an expiration time.
type session struct {
	ID     string
	expiry time.Time
}

// isExpired returns true if the session has expired,
// false otherwise.
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// userOrDefault returns an authenticated User
// from the Authenticator session store, or an empty one.
func (a *Authenticator) userOrDefault(id string) User {
	session, ok := a.sessions[id]
	if !ok {
		return User{}
	}
	return User{ID: session.ID, Authenticated: true}
}

// Identify retrieves a session or creates a new one.
func (a *Authenticator) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")

		var sessionToken string

		switch err {
		case nil:
			sessionToken = c.Value

		case http.ErrNoCookie:
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(a.maxSessionLength)

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := a.userOrDefault(sessionToken)

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
