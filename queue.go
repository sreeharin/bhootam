package bhootam

import "github.com/google/uuid"

type Queue struct {
	jobs chan Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan Job)}
}

func (q *Queue) AddTask(task Task) string {
	id := uuid.NewString()

	q.jobs <- Job{id, task}

	return id
}
