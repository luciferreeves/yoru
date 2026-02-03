package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"yoru/models"
	"yoru/repository"
	"yoru/shared"
	"yoru/types"

	"golang.org/x/crypto/ssh"
)

// NewClient creates a new SSH client instance
func NewClient(host *models.Host, credential any) *client {
	return &client{
		host:       host,
		credential: credential,
		state:      stateConnecting,
		outputChan: make(chan []byte, 100),
		errorChan:  make(chan error, 10),
	}
}

// Connect establishes an SSH connection
func (c *client) Connect() error {
	// Build SSH configuration
	config, err := BuildSSHConfig(c.credential)
	if err != nil {
		return fmt.Errorf("failed to build SSH config: %w", err)
	}

	// Connect to SSH server
	addr := net.JoinHostPort(c.host.Hostname, fmt.Sprintf("%d", c.host.Port))

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Starting connection to %s port %d", c.host.Hostname, c.host.Port),
	})

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Starting address resolution of %s", c.host.Hostname),
	})

	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: "- Address resolution finished",
	})

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Connecting to %s port %d", c.host.Hostname, c.host.Port),
	})

	// Custom host key callback to capture and verify the key
	var serverKey ssh.PublicKey
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		serverKey = key
		return nil
	}

	// Establish SSH connection
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Connection to %s established", c.host.Hostname),
	})

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: "- Starting SSH session",
	})

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Remote server: %s", string(sshConn.ServerVersion())),
	})

	// Verify host key
	knownHost, err := VerifyHostKey(c.host.Hostname, c.host.Port, serverKey)
	if err != nil {
		// Host key not known — ask user whether to add to known hosts
		fingerprint := GetFingerprint(serverKey)
		c.hostKeyDecision = make(chan bool, 1)

		shared.SendMessage(types.SSHHostKeyMsg{
			HostID:      c.host.ID,
			Hostname:    c.host.Hostname,
			Port:        c.host.Port,
			KeyType:     serverKey.Type(),
			Fingerprint: fingerprint,
			ServerKey:   serverKey,
		})

		// Block until user decides
		if save := <-c.hostKeyDecision; save {
			_ = SaveHostKey(c.host.Hostname, c.host.Port, serverKey)
			shared.SendMessage(types.SSHConnectingMsg{
				HostID:  c.host.ID,
				Message: "- Host key added to known hosts",
			})
		} else {
			shared.SendMessage(types.SSHConnectingMsg{
				HostID:  c.host.ID,
				Message: "- Host key not saved",
			})
		}
	} else {
		shared.SendMessage(types.SSHConnectingMsg{
			HostID:  c.host.ID,
			Message: fmt.Sprintf("- Checking host key: %s", knownHost.Fingerprint),
		})

		shared.SendMessage(types.SSHConnectingMsg{
			HostID:  c.host.ID,
			Message: fmt.Sprintf("- Host %s:%d is known and matches", c.host.Hostname, c.host.Port),
		})
	}

	// Create SSH client
	c.sshClient = ssh.NewClient(sshConn, chans, reqs)

	shared.SendMessage(types.SSHAuthenticatingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Authenticating to %s:%d", c.host.Hostname, c.host.Port),
	})

	// Determine auth method
	authMethod := "unknown"
	if _, ok := c.credential.(*models.Identity); ok {
		authMethod = "password"
	} else if _, ok := c.credential.(*models.Key); ok {
		authMethod = "publickey"
	}

	shared.SendMessage(types.SSHAuthenticatingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Authenticating using %s method", authMethod),
	})

	shared.SendMessage(types.SSHAuthenticatingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Authentication succeeded (%s)", authMethod),
	})

	shared.SendMessage(types.SSHAuthenticatingMsg{
		HostID:  c.host.ID,
		Message: fmt.Sprintf("- Authenticated to %s:%d", c.host.Hostname, c.host.Port),
	})

	return nil
}

// StartSession creates and starts an SSH session with a PTY
func (c *client) StartSession(width, height int) error {
	if c.sshClient == nil {
		return fmt.Errorf("SSH client not connected")
	}

	c.termWidth = width
	c.termHeight = height

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: "- Creating terminal session",
	})

	// Create new session
	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	c.session = session

	// Request PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("failed to request PTY: %w", err)
	}

	// Setup I/O pipes
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	c.stdin = stdin

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	shared.SendMessage(types.SSHConnectingMsg{
		HostID:  c.host.ID,
		Message: "- Shell started successfully",
	})

	c.state = stateConnected

	// Start output streaming
	go c.streamOutput(stdout, stderr)

	// Create connection log
	localHostname, _ := os.Hostname()
	localIP := getLocalIP()

	connectionLog := &models.ConnectionLog{
		StartedAt:      time.Now(),
		LocalHostname:  localHostname,
		LocalIP:        localIP,
		RemoteHostname: c.host.Hostname,
		Mode:           c.host.Mode,
		CredentialID:   c.host.CredentialID,
		CredentialType: c.host.CredentialType,
	}

	// Save to database
	if err := repository.CreateConnectionLog(connectionLog); err == nil {
		c.connectionLog = connectionLog
	}

	// Send connected message
	shared.SendMessage(types.SSHConnectedMsg{
		HostID:        c.host.ID,
		Client:        c,
		ConnectionLog: connectionLog,
	})

	return nil
}

// streamOutput reads from stdout/stderr and sends to program
func (c *client) streamOutput(stdout, stderr io.Reader) {
	// Merge stdout and stderr
	reader := io.MultiReader(stdout, stderr)
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])

			shared.SendMessage(types.SSHOutputMsg{
				HostID: c.host.ID,
				Data:   data,
			})
		}

		if err != nil {
			if err == io.EOF {
				shared.SendMessage(types.SSHDisconnectedMsg{
					HostID: c.host.ID,
				})
			} else {
				shared.SendMessage(types.SSHErrorMsg{
					HostID: c.host.ID,
					Error:  fmt.Errorf("output stream error: %w", err),
				})
			}
			break
		}
	}

	c.state = stateDisconnected
}

// SendInput sends input to the SSH session
func (c *client) SendInput(data []byte) error {
	if c.stdin == nil {
		return fmt.Errorf("stdin not available")
	}

	_, err := c.stdin.Write(data)
	return err
}

// Resize resizes the terminal
func (c *client) Resize(width, height int) error {
	if c.session == nil {
		return fmt.Errorf("session not available")
	}

	c.termWidth = width
	c.termHeight = height

	return c.session.WindowChange(height, width)
}

// Close closes the SSH connection
func (c *client) Close() error {
	c.state = stateDisconnected

	// Update connection log
	if c.connectionLog != nil {
		endedAt := time.Now()
		c.connectionLog.EndedAt = &endedAt
		repository.UpdateConnectionLog(c.connectionLog)
	}

	// Close session
	if c.session != nil {
		c.session.Close()
	}

	// Close SSH client
	if c.sshClient != nil {
		c.sshClient.Close()
	}

	return nil
}

// getLocalIP gets the local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "unknown"
}