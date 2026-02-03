package ssh

import (
	"fmt"
	"yoru/models"
	"yoru/repository"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

// activeClients stores active SSH clients by host ID
var activeClients = make(map[uint]*client)

// InitiateConnection starts an SSH connection asynchronously
func InitiateConnection(host *models.Host) tea.Cmd {
	return func() tea.Msg {
		// Start connection in goroutine
		go connectAsync(host)

		// Return connecting message immediately
		return types.SSHConnectingMsg{
			HostID:  host.ID,
			Message: "- Initializing connection",
		}
	}
}

// connectAsync performs the full connection flow
func connectAsync(host *models.Host) {
	// Load credential
	credential, err := LoadCredential(host)
	if err != nil {
		shared.SendMessage(types.SSHErrorMsg{
			HostID: host.ID,
			Error:  fmt.Errorf("failed to load credential: %w", err),
		})
		return
	}

	// Create client
	client := NewClient(host, credential)

	// Store active client
	activeClients[host.ID] = client

	// Attempt connection (blocks on host key decision if key is unknown)
	if err := client.Connect(); err != nil {
		shared.SendMessage(types.SSHErrorMsg{
			HostID: host.ID,
			Error:  err,
		})
		return
	}

	// Start session with dynamic dimensions
	// Note: Dimensions will be updated by terminal screen after creation
	// Using default 80x24 initially, will be resized immediately
	if err := client.StartSession(80, 24); err != nil {
		shared.SendMessage(types.SSHErrorMsg{
			HostID: host.ID,
			Error:  fmt.Errorf("failed to start session: %w", err),
		})
		return
	}
}

// ContinueAfterHostKeyVerification unblocks the connection goroutine after the user
// decides whether to save the host key. save=true adds it to known hosts.
func ContinueAfterHostKeyVerification(hostID uint, save bool) {
	client, ok := activeClients[hostID]
	if !ok {
		shared.SendMessage(types.SSHErrorMsg{
			HostID: hostID,
			Error:  fmt.Errorf("client not found"),
		})
		return
	}

	client.hostKeyDecision <- save
}

// RetryConnection retries a failed connection
func RetryConnection(hostID uint) {
	if client, ok := activeClients[hostID]; ok {
		// Close existing client
		client.Close()
	}

	// Get host from repository and retry
	host, err := repository.GetHostByID(hostID)
	if err != nil {
		shared.SendMessage(types.SSHErrorMsg{
			HostID: hostID,
			Error:  fmt.Errorf("failed to get host: %w", err),
		})
		return
	}

	// Retry connection
	go connectAsync(host)
}

// GetClient returns the active client for a host
func GetClient(hostID uint) *client {
	return activeClients[hostID]
}

// CloseConnection closes an active SSH connection
func CloseConnection(hostID uint) {
	if client, ok := activeClients[hostID]; ok {
		client.Close()
		delete(activeClients, hostID)
	}
}

// ResizeTerminal resizes the terminal for an active connection
func ResizeTerminal(hostID uint, width, height int) error {
	client, ok := activeClients[hostID]
	if !ok {
		return fmt.Errorf("client not found")
	}

	return client.Resize(width, height)
}

// SendInput sends keyboard input to an active connection
func SendInput(hostID uint, data []byte) error {
	client, ok := activeClients[hostID]
	if !ok {
		return fmt.Errorf("client not found")
	}

	return client.SendInput(data)
}