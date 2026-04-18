package bhootam

type Args []any
type Func func(Args) Outcome

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
}

// Run executes the function with the provided arguments
func (t Task) Run() Outcome {
	return t.Function(t.Args)
}
