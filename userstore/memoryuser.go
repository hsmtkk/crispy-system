package userstore

import "context"

type MemoryImpl struct {
	userVisitedMap map[string]int
}

func NewMemoryImpl() UserStore {
	userVisitedMap := map[string]int{}
	return &MemoryImpl{userVisitedMap}
}

func (m *MemoryImpl) Increment(ctx context.Context, userID string) (int, error) {
	count, ok := m.userVisitedMap[userID]
	newCount := 1
	if ok {
		newCount = count + 1
	}
	m.userVisitedMap[userID] = newCount
	return count, nil
}
