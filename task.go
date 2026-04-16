package bhootam

type Args []any
type Func func(Args) Value

type Value struct {
	Value any
	Err   error
}

type Task struct {
	Function Func
	Args     Args
}

func (t Task) Run() Value {
	return t.Function(t.Args)
}
