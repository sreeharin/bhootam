package bhootam

import (
	"errors"
	"sync"
)

type Result struct {
	Outcome
	Status JobState
}

type Store struct {
	data map[string]Result
	mu   sync.Mutex
}

func NewStore() *Store {
	return &Store{data: make(map[string]Result)}
}

func (s *Store) Set(key string, value Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Store) Get(key string) (Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]

	if !ok {
		return Result{}, errors.New("Job id not found!")
	}

	return val, nil
}
