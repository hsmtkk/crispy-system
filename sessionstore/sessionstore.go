package sessionstore

type SessionStore interface {
	NewSession(userID string) (string, error)
	GetUserID(sessionID string) (string, error)
	DeleteSession(sessionID string) error
}
