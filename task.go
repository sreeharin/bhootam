package bhootam

type Args []any

type Value struct {
	Value any
	Err   error
}

type Task struct {
	Function func(Args) Value
	Args     Args
}

func (t Task) Run() Value {
	return t.Function(t.Args)
}
