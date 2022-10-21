package sessionstore

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type MemoryImpl struct {
	sessionUserMap map[string]string
}

func NewMemoryImpl() SessionStore {
	return &MemoryImpl{sessionUserMap: map[string]string{}}
}

func (m *MemoryImpl) NewSession(userID string) (string, error) {
	sessionID := uuid.NewString()
	m.sessionUserMap[sessionID] = userID
	log.Printf("NewSession: %s %v", userID, m.sessionUserMap)
	return sessionID, nil
}

func (m *MemoryImpl) GetUserID(sessionID string) (string, error) {
	userID, ok := m.sessionUserMap[sessionID]
	if ok {
		log.Printf("GetUserID: %s %s", sessionID, userID)
		return userID, nil
	} else {
		return "", fmt.Errorf("session %s does not exist", sessionID)
	}
}

func (m *MemoryImpl) DeleteSession(sessionID string) error {
	delete(m.sessionUserMap, sessionID)
	log.Printf("DeleteSession: %s %v", sessionID, m.sessionUserMap)
	return nil
}
