package bhootam

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	jobs chan Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan Job)}
}

// AddTask enqueues Task to the queue (Job channel)
// a UUID is generated per unique job
func (q *Queue) AddTask(task *Task) (string, <-chan struct{}, <-chan struct{}) {
	id := uuid.NewString()
	ack := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	DEFAULT_TIMEOUT := 1 * time.Minute

	// Set default timeout if no timeout was set
	if task.Timeout == 0 {
		task.Timeout = DEFAULT_TIMEOUT
	}

	ctx, cancel := context.WithTimeout(context.TODO(), task.Timeout)
	q.jobs <- Job{ctx, cancel, id, task, ack, done}

	return id, ack, done
}
