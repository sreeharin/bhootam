package bhootam

import "testing"

func sum(args Args) Value {
	ans := args[0].(int) + args[1].(int)
	return Value{Value: ans, Err: nil}
}

func TestRun(t *testing.T) {
	task := Task{Function: sum, Args: Args{5, 5}}
	ret := task.Run()

	if ret.Value.(int) != 10 {
		t.Error("Incorrect value")
	}
}
