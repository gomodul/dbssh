package dbssh

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSH ...
type SSH struct {
	Net net.Conn
	SSH *ssh.Client
}

// Open ...
func Open(c Config, driver Driver) (*SSH, string, error) {
	var (
		err         error
		agentClient agent.Agent

		conn = new(SSH)
	)

	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	// Establish a connection to the local ssh-agent
	if conn.Net, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn.Net)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	// When there's a non-empty password add the password AuthMethod
	sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
		return c.Pass, nil
	}))

	// Connect to the SSH Server
	if conn.SSH, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), sshConfig); err != nil {
		return nil, "", err
	}

	driver.Register(conn.SSH)
	return conn, driver.Name(), nil
}

// Close ...
func (conn SSH) Close() {
	if conn.SSH != nil {
		_ = conn.SSH.Close()
	}
	if conn.Net != nil {
		_ = conn.Net.Close()
	}
}
