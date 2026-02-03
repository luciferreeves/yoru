package types

import (
	"golang.org/x/crypto/ssh"
)

// SSH Bubble Tea messages for async events

type SSHConnectingMsg struct {
	HostID  uint
	Message string
}

type SSHAuthenticatingMsg struct {
	HostID  uint
	Message string
}

type SSHHostKeyMsg struct {
	HostID      uint
	Hostname    string
	Port        int
	KeyType     string
	Fingerprint string
	ServerKey   ssh.PublicKey
}

type SSHConnectedMsg struct {
	HostID        uint
	Client        any // *ssh.Client from ssh package
	ConnectionLog any // *models.ConnectionLog
}

type SSHOutputMsg struct {
	HostID uint
	Data   []byte
}

type SSHErrorMsg struct {
	HostID uint
	Error  error
}

type SSHDisconnectedMsg struct {
	HostID uint
}