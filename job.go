package bhootam

import (
	"context"
)

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
