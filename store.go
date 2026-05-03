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
	// store the return value from function
	data map[string]Result

	// to keep track of the jobs
	// to make sure a job is run only once at a given time
	jobs map[string]bool

	mu sync.Mutex
}

func NewStore() *Store {
	return &Store{data: make(map[string]Result), jobs: make(map[string]bool)}
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

func (s *Store) Acquire(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// job is being processed by a worker
	if s.jobs[id] {
		return false
	}

	s.jobs[id] = true
	return true
}

func (s *Store) Release(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.jobs, id)
}
