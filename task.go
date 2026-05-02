package bhootam

import "time"

type Args []any
type Func func(Args) Outcome
type option func() func(*Task)

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
	Retry    int
}

// Run executes the function with the provided arguments
func (t Task) Run() Outcome {
	return t.Function(t.Args)
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

func withArgs(args Args) func(*Task) {
	return func(t *Task) {
		t.Args = args
	}
}

func withTimeout(timeout time.Duration) func(*Task) {
	return func(t *Task) {
		t.Timeout = timeout
	}
}

func withRetries(retry int) func(*Task) {
	return func(t *Task) {
		if retry > 0 {
			t.Retry = retry
		}
	}
}
