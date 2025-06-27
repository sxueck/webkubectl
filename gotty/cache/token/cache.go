package token

import "time"

// TtyParameter kubectl tty param
type TtyParameter struct {
	Title string
	Arg   string
}

// SessionInfo represents an active session
type SessionInfo struct {
	SessionToken string
	WorkingDir   string
	CreatedAt    int64
	LastAccess   int64
}

// interface that defines token cache behavior
type Cache interface {
	Get(token string) *TtyParameter
	Delete(token string) error
	Add(token string, param *TtyParameter, d time.Duration) error

	// Session management
	GetSession(sessionToken string) *SessionInfo
	SetSession(sessionToken string, session *SessionInfo, d time.Duration) error
	UpdateSessionAccess(sessionToken string) error
	DeleteSession(sessionToken string) error
}
