package shared

import tea "github.com/charmbracelet/bubbletea"

// Program is a reference to the Bubble Tea program for sending messages from goroutines
var Program *tea.Program

// SetProgram sets the global program reference
func SetProgram(p *tea.Program) {
	Program = p
}

// SendMessage sends a message to the program
func SendMessage(msg tea.Msg) {
	if Program != nil {
		Program.Send(msg)
	}
}