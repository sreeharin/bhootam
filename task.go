package bhootam

type Args []any
type Func func(Args) Outcome

type Outcome struct {
	Value any
	Err   error
}

type Task struct {
	Function Func
	Args     Args
}

func (t Task) Run() Outcome {
	return t.Function(t.Args)
}
