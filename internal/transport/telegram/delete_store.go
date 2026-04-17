package telegram

import "sync"

type pendingDelete struct {
	userID    int64
	plantID   int64
	plantName string
}

type pendingDeleteKey struct {
	telegramUserID int64
	messageID      int
}

type PendingDeleteStore struct {
	mu    sync.RWMutex
	items map[pendingDeleteKey]pendingDelete
}

func NewPendingDeleteStore() *PendingDeleteStore {
	return &PendingDeleteStore{
		items: make(map[pendingDeleteKey]pendingDelete),
	}
}

func (s *PendingDeleteStore) Set(telegramUserID int64, messageID int, value pendingDelete) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[pendingDeleteKey{telegramUserID: telegramUserID, messageID: messageID}] = value
}

func (s *PendingDeleteStore) Get(telegramUserID int64, messageID int) (pendingDelete, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.items[pendingDeleteKey{telegramUserID: telegramUserID, messageID: messageID}]
	return value, ok
}

func (s *PendingDeleteStore) Clear(telegramUserID int64, messageID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, pendingDeleteKey{telegramUserID: telegramUserID, messageID: messageID})
}

func (s *PendingDeleteStore) ClearAllForUser(telegramUserID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.items {
		if key.telegramUserID == telegramUserID {
			delete(s.items, key)
		}
	}
}
