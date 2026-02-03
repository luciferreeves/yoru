package ssh

import (
	"io"
	"yoru/models"

	"golang.org/x/crypto/ssh"
)

// connection states
type connectionState int

const (
	stateConnecting connectionState = iota
	stateAuthenticating
	stateVerifyingHost
	stateConnected
	stateDisconnected
	stateError
)

// client is the SSH client wrapper
type client struct {
	host           *models.Host
	credential     any // *models.Identity or *models.Key
	sshClient      *ssh.Client
	session        *ssh.Session
	state          connectionState
	connectionLog  *models.ConnectionLog

	// terminal dimensions
	termWidth      int
	termHeight     int

	// i/o
	stdin          io.WriteCloser
	outputChan     chan []byte
	errorChan      chan error

	// host key verification: receives true to save, false to skip
	hostKeyDecision chan bool
}