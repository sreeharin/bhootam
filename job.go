package bhootam

import (
	"context"
)

type jobOption func(*Job)

type JobState string

const (
	JobCompleted JobState = "completed"
	JobRunning   JobState = "running"
	JobError     JobState = "error"
	JobTimeOut   JobState = "timeout"
	JobRetry     JobState = "retry"
)

type Job struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	id        string
	task      *Task

	// ack is the acknowledgment provided by the worker
	// to inform that a worker has picken up the task
	ack chan struct{}

	// done is used to infrom that the task has completed running
	done chan struct{}

	// mu sync.Mutex
}

func NewJob(ctx context.Context, ctxCancel context.CancelFunc, id string, task *Task, options ...jobOption) *Job {
	job := &Job{
		ctx:       ctx,
		ctxCancel: ctxCancel,
		id:        id,
		task:      task,
	}

	for _, opt := range options {
		opt(job)
	}

	return job
}

func withAck(ack chan struct{}) jobOption {
	return func(j *Job) {
		j.ack = ack
	}
}

func withDone(done chan struct{}) jobOption {
	return func(j *Job) {
		j.done = done
	}
}
