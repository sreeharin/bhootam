package bhootam

import (
	"sync/atomic"
	"time"
)

type Args []any
type Func func(Args) Outcome
type option func(*Task)

// Outcome is the return value we expect from a function
type Outcome struct {
	Value any
	Err   error
}

// Task is an abstraction layer
// We expect the user functions to be wrapped inside Task
type Task struct {
	Function Func
	Args     Args
	Timeout  time.Duration
	Retry    atomic.Int32
}

// Run executes the function with the provided arguments
func (t *Task) Run() Outcome {
	// res := make(chan Outcome, 1)
	// res <- t.Function(t.Args)
	return t.Function(t.Args)
	// return res
}

func (t *Task) DecrementRetry() {
	t.Retry.Add(-1)
}

func NewTask(function Func, options ...option) *Task {
	task := &Task{
		Function: function,
	}

	for _, opt := range options {
		opt(task)
	}

	return task
}

func withArgs(args Args) option {
	return func(t *Task) {
		t.Args = args
	}
}

func withTimeout(timeout time.Duration) option {
	return func(t *Task) {
		t.Timeout = timeout
	}
}

func withRetry(count int32) option {
	return func(t *Task) {
		if count > 0 {
			t.Retry.Add(count)
		}
	}
}
