package sessionstore

import (
	"context"
	"fmt"
	"log"
)

type MemoryImpl struct {
	sessionUserMap map[string]string
}

func NewMemoryImpl() SessionStore {
	return &MemoryImpl{sessionUserMap: map[string]string{}}
}

func (m *MemoryImpl) NewSession(ctx context.Context, sessionID, userID string) (string, error) {
	m.sessionUserMap[sessionID] = userID
	log.Printf("NewSession: %s %v", userID, m.sessionUserMap)
	return sessionID, nil
}

func (m *MemoryImpl) GetUserID(ctx context.Context, sessionID string) (string, error) {
	userID, ok := m.sessionUserMap[sessionID]
	if ok {
		log.Printf("GetUserID: %s %s", sessionID, userID)
		return userID, nil
	} else {
		return "", fmt.Errorf("session %s does not exist", sessionID)
	}
}

func (m *MemoryImpl) DeleteSession(ctx context.Context, sessionID string) error {
	delete(m.sessionUserMap, sessionID)
	log.Printf("DeleteSession: %s %v", sessionID, m.sessionUserMap)
	return nil
}
