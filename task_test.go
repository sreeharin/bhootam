package bhootam

import "testing"

func sampleSumTask(args Args) Outcome {
	ans := args[0].(int) + args[1].(int)
	return Outcome{Value: ans, Err: nil}
}

func TestRun(t *testing.T) {
	task := Task{function: sampleSumTask, args: Args{5, 5}}
	res := task.Run()
	// res := <-outcome

	if res.Value.(int) != 10 {
		t.Error("Incorrect value")
	}
}
