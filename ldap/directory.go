package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap"
)

// A directory represents an LDAP search domain.
type directory struct {
	bindAddr   url.URL
	bindUser   string
	bindPass   string
	baseDN     string
	classValue string
	idKey      string
}

// Bind is a functional option that verifies both the connection
// and bind status to a directory with the credentials provided to it.
func Bind(username string, password string) func(d *directory) error {
	return func(d *directory) error {
		conn, err := ldap.DialURL(d.bindAddr.String())
		if err != nil {
			return err
		}

		defer conn.Close()

		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}

		d.bindUser = username
		d.bindPass = password

		return conn.Bind(d.bindUser, d.bindPass)
	}
}

// Fields is a functional option that sets parameters for the LDAP filter.
// It sets both the value for the parameter "objectClass"
// and the key name for the username
func Fields(classValue string, idKey string) func(d *directory) {
	return func(d *directory) {
		d.classValue = classValue
		d.idKey = idKey
	}
}

// NewDirectory initializes an ldap client. The initialization fails if any
// functional option returns an error.
func NewDirectory(addr string, baseDN string, options ...func(*directory) error) (*directory, error) {

	url, err := url.Parse(addr)
	if err != nil {
		return nil, errors.New(addr + " is an invalid URL: " + err.Error())
	}

	d := directory{
		bindAddr:   *url,
		baseDN:     baseDN,
		classValue: "organizationalPerson",
		idKey:      "uid",
	}

	for _, option := range options {
		err := option(&d)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

// Authenticate returns an error if the username is not found within
// the directory or the username does not bind to it with the provided password.
func (d directory) Authenticate(username string, password string) error {

	conn, err := ldap.DialURL(d.bindAddr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	err = conn.Bind(d.bindUser, d.bindPass)
	if err != nil {
		return err
	}

	searchRequest := ldap.NewSearchRequest(
		d.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=%s)(%s=%s))", d.classValue, d.idKey, ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("user %q does not exist or too many entries found", username)
	}

	userdn := sr.Entries[0].DN

	err = conn.Bind(userdn, password)
	if err != nil {
		return err
	}

	return nil
}
