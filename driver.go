package dbssh

import (
	"golang.org/x/crypto/ssh"
)

// Driver ...
type Driver interface {
	Name() string
	Register(*ssh.Client)
}
