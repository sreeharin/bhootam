package bhootam

import (
	"sync/atomic"
	"time"
)

type Args []any
type Func func(Args) Outcome
type taskOption func(*Task)

// Outcome is the return value we expect from a function
type Outcome struct {
	Value any
	Err   error
}

// Task is an abstraction layer
// We expect the user functions to be wrapped inside Task
type Task struct {
	function Func
	args     Args
	timeout  time.Duration
	retry    atomic.Int32
}

// Run executes the function with the provided arguments
func (t *Task) Run() Outcome {
	return t.function(t.args)
}

// DecrementRetry reduces the retry count
// if the number is 0, we stop retrying
func (t *Task) DecrementRetry() {
	if t.retry.Load() > 0 {
		t.retry.Add(-1)
	}
}

func NewTask(function Func, options ...taskOption) *Task {
	task := &Task{
		function: function,
	}

	for _, opt := range options {
		opt(task)
	}

	return task
}

func withArgs(args Args) taskOption {
	return func(t *Task) {
		t.args = args
	}
}

func withTimeout(timeout time.Duration) taskOption {
	return func(t *Task) {
		t.timeout = timeout
	}
}

func withRetry(count int32) taskOption {
	return func(t *Task) {
		if count > 0 {
			t.retry.Add(count)
		}
	}
}
