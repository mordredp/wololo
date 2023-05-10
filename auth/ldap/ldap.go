package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap"
)

type Authenticator struct {
	bindAddr url.URL
	bindUser string
	bindPass string
	baseDN   string
}

func New(addr string, username, password string, baseDN string) (*Authenticator, error) {

	url, err := url.Parse(addr)
	if err != nil {
		return nil, errors.New(addr + " is an invalid URL: " + err.Error())
	}

	authenticator := new(Authenticator)

	authenticator.bindAddr = *url
	authenticator.bindUser = username
	authenticator.bindPass = password

	return authenticator, nil
}

func (a Authenticator) MustAuthenticate(username string, password string) error {

	l, err := ldap.DialURL(a.bindAddr.String())
	if err != nil {
		return err
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	// First bind with a read only user
	err = l.Bind(a.bindUser, a.bindPass)
	if err != nil {
		return err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		a.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return errors.New("user does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}

	return nil
}
