package bhootam

import "github.com/google/uuid"

type Queue struct {
	jobs chan Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan Job)}
}

// AddTask enqueues Task to the queue (Job channel)
// a UUID is generated per unique job
func (q *Queue) AddTask(task Task) (string, <-chan struct{}) {
	id := uuid.NewString()
	ack := make(chan struct{}, 1)

	q.jobs <- Job{id, task, ack}

	return id, ack
}
