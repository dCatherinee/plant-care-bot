package telegram

import "sync"

type State string

const (
	StateIdle             State = "idle"
	StateWaitingPlantName State = "waiting_plant_name"
)

type StateStore struct {
	mu     sync.RWMutex
	states map[int64]State
}

func NewStateStore() *StateStore {
	return &StateStore{
		states: make(map[int64]State),
	}
}

func (s *StateStore) Set(userID int64, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[userID] = state
}

func (s *StateStore) Get(userID int64) State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.states[userID]
	if !ok {
		return StateIdle
	}

	return state
}

func (s *StateStore) Clear(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.states, userID)
}
