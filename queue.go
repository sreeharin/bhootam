package bhootam

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	jobs  chan *Job
	count atomic.Int32
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan *Job)}
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
	q.jobs <- &Job{ctx: ctx, ctxCancel: cancel, id: id, task: task, ack: ack, done: done}

	q.count.Add(1)

	return id, ack, done
}

// func (q *Queue) Enqueue (job *Job) {
// 	if
// }
