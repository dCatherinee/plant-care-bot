package telegram

import "sync"

type pendingDelete struct {
	userID    int64
	plantID   int64
	plantName string
}

type PendingDeleteStore struct {
	mu    sync.RWMutex
	items map[int64]pendingDelete
}

func NewPendingDeleteStore() *PendingDeleteStore {
	return &PendingDeleteStore{
		items: make(map[int64]pendingDelete),
	}
}

func (s *PendingDeleteStore) Set(telegramUserID int64, value pendingDelete) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[telegramUserID] = value
}

func (s *PendingDeleteStore) Get(telegramUserID int64) (pendingDelete, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.items[telegramUserID]
	return value, ok
}

func (s *PendingDeleteStore) Clear(telegramUserID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, telegramUserID)
}
