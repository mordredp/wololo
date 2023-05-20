package auth

import (
	"log"
	"text/template"
	"time"

	"github.com/mordredp/wololo/provider"
	"github.com/mordredp/wololo/provider/ldap"
)

// Authenticator manages sessions and authentication providers.
type Authenticator struct {
	sessions         map[string]session
	maxSessionLength time.Duration
	lastCleanup      time.Time
	tpl              *template.Template
	providers        []provider.Provider
}

// LDAP is a functional option that instantiates an LDAP provider
// for the Authenticator.
func LDAP(addr string, baseDN string, username string, password string) func(a *Authenticator) error {
	return func(a *Authenticator) error {

		ldap, err := ldap.NewDirectory(
			addr,
			baseDN,
			ldap.Bind(username, password))

		if err != nil {
			return err
		}

		a.providers = append(a.providers, ldap)
		log.Printf("configured LDAP provider on %q with base DN %q", addr, baseDN)

		return nil
	}
}

// Static is a functional option that instantiates a Static provider
// for the Authenticator.
func Static(password string) func(a *Authenticator) error {
	return func(a *Authenticator) error {
		a.providers = append(a.providers, provider.Static(password))
		log.Printf("configured Static provider ")

		return nil
	}
}

// New initializes a new Authenticator.
func New(sessionSeconds int, options ...func(*Authenticator) error) *Authenticator {

	a := Authenticator{
		sessions:         make(map[string]session),
		maxSessionLength: time.Duration(sessionSeconds) * time.Second,
		lastCleanup:      time.Now(),
		tpl:              template.Must(template.ParseGlob("auth/templates/*.gohtml")),
		providers:        make([]provider.Provider, 0),
	}

	for _, option := range options {
		err := option(&a)
		if err != nil {
			log.Printf("options: %s", err)
			continue
		}
	}

	return &a
}
