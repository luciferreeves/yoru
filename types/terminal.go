package types

// KeyCaptureMode defines how keyboard input is handled
type KeyCaptureMode int

const (
	KeyCaptureNormal   KeyCaptureMode = iota // Manager handles keys
	KeyCaptureTerminal                        // Terminal captures all keys
)

// KeyCapturer is implemented by screens that support key capture mode
type KeyCapturer interface {
	GetKeyCaptureMode() KeyCaptureMode
}