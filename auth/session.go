package auth

import (
	"text/template"
	"time"
)

type key int

const (
	// UserKey is the key to the value of a User in a context
	UserKey key = iota
	// maxSessionLength is the maximum session duration (in seconds)
	maxSessionLength time.Duration = time.Duration(120) * time.Second
)

// User holds a users account information
type User struct {
	ID            string
	Authenticated bool
}

// Credentials holds the username and password used
// to authenticate a session
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// session contains the username of the user it is assigned to
// and its expiration time
type session struct {
	ID     string
	expiry time.Time
}

// isExpired returns true if the session has expired,
// false otherwise
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

var sessions = map[string]session{}

// tpl holds templates data
var tpl *template.Template

var lastCleanup time.Time

func init() {
	tpl = template.Must(template.ParseGlob("auth/templates/*.gohtml"))
	lastCleanup = time.Now()
}

// userOrNew returns a user from session s or an empty user
func userOrNew(id string) User {
	session, ok := sessions[id]
	if !ok {
		return User{Authenticated: false}
	}
	return User{ID: session.ID, Authenticated: true}
}
