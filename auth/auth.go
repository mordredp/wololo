package auth

import (
	"log"
	"text/template"
	"time"

	"github.com/mordredp/wololo/ldap"
)

// Authenticator manages sessions and authentication providers.
type Authenticator struct {
	sessions         map[string]session
	maxSessionLength time.Duration
	lastCleanup      time.Time
	tpl              *template.Template
	provider         Provider
}

// LDAP is a functional option that instantiates an LDAP provider
// for the Authenticator
func LDAP(addr string, baseDN string, username string, password string) func(a *Authenticator) error {
	return func(a *Authenticator) error {

		ldap, err := ldap.NewDirectory(
			addr,
			baseDN,
			ldap.Bind(username, password))

		if err != nil {
			return err
		}

		a.provider = ldap
		log.Printf("LDAP provider on %q with base DN %q", addr, baseDN)

		return nil
	}
}

// New initializes a new Authenticator. The initialization fails if any
// functional option returns an error.
func New(sessionSeconds int, options ...func(*Authenticator) error) (*Authenticator, error) {

	a := Authenticator{
		sessions:         make(map[string]session),
		maxSessionLength: time.Duration(sessionSeconds) * time.Second,
		lastCleanup:      time.Now(),
		tpl:              template.Must(template.ParseGlob("auth/templates/*.gohtml")),
		provider:         &nullProvider{},
	}

	for _, option := range options {
		err := option(&a)
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}
