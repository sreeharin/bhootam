package bhootam

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const DefaultTimeout = 1 * time.Minute

type Queue struct {
	jobs  chan *Job
	count atomic.Int32
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan *Job)}
}

// NewTask enqueues a new Task to the queue (Job channel)
// a UUID is generated per unique job
func (q *Queue) CreateJob(task *Task) (string, <-chan struct{}, <-chan struct{}) {
	id := uuid.NewString()
	ack := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	// Set default timeout if no timeout was set
	if task.Timeout == 0 {
		task.Timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(context.TODO(), task.Timeout)

	job := NewJob(ctx, cancel, id, task, withAck(ack), withDone(done))
	q.Enqueue(job)

	return id, ack, done
}

func (q *Queue) Enqueue(job *Job) {
	q.jobs <- job
	q.count.Add(1)
}
