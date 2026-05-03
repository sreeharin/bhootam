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

// CreateJob creates a new job from the task
// creates a UUID for the job with ack, and done channels
// sets the timeout if there is none
func (q *Queue) CreateJob(task *Task) (string, <-chan struct{}, <-chan struct{}) {
	id := uuid.NewString()
	ack := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	// Set default timeout if no timeout was set
	if task.timeout == 0 {
		task.timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(context.TODO(), task.timeout)

	job := NewJob(ctx, cancel, id, task, withAck(ack), withDone(done))
	q.Enqueue(job)

	return id, ack, done
}

func (q *Queue) Enqueue(job *Job) {
	q.jobs <- job
	q.count.Add(1)
}
