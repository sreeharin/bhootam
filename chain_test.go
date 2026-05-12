package bhootam

import (
	"testing"
)

func TestChain(t *testing.T) {
	sum := func(args Args) Outcome {
		ans := args[0].(int) + args[1].(int)
		return Outcome{Value: ans, Err: nil}
	}

	mul := func(arg Args) Outcome {
		ans := arg[0].(int) * arg[1].(int)
		return Outcome{Value: ans, Err: nil}
	}

	task1 := NewTask(sum, WithArgs(Args{5, 5}))
	task2 := NewTask(mul, WithArgs(Args{10}))

	Chain(task1, task2)
	id, _, done := Chain(task1, task2)
	<-done

	if res, err := store.Get(id); err != nil {
		t.Fatalf("failed to get result from store: %v", err)
	} else {
		if res.Status != JobCompleted {
			t.Errorf("Wrong status. Expected: %s, Got: %s", JobCompleted, res.Status)
		}

		if res.Value != 100 {
			t.Errorf("Wrong value. Expected: %v Got: %v", 100, res.Value)
		}
	}
}
