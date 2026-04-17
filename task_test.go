package bhootam

import "testing"

func sampleSumTask(args Args) Outcome {
	ans := args[0].(int) + args[1].(int)
	return Outcome{Value: ans, Err: nil}
}

func TestRun(t *testing.T) {
	task := Task{Function: sampleSumTask, Args: Args{5, 5}}
	ret := task.Run()

	if ret.Value.(int) != 10 {
		t.Error("Incorrect value")
	}
}
