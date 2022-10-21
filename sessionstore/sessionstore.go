package sessionstore

import "context"

type SessionStore interface {
	NewSession(ctx context.Context, sessionID, userID string) (string, error)
	GetUserID(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
