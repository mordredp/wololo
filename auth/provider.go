package auth

import "fmt"

// A Provider can Authenticate a pair of username and password.
type Provider interface {
	// Authenticate returns an error if the username and password are not valid.
	Authenticate(username string, password string) error
}

// nullProvider is a default implementation of a Provider.
type nullProvider struct{}

// Authenticate always returns an error signaling that
// the default provider is being used.
func (n *nullProvider) Authenticate(username string, password string) error {
	return fmt.Errorf("nullProvider: no provider defined")
}
