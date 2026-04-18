package bhootam

import (
	"os"
	"testing"
)

var (
	q     *Queue
	store *Store
)

func TestMain(m *testing.M) {
	q = NewQueue()
	store = NewStore()

	StartWorker(q, store)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestStartWorker(t *testing.T) {
	task := Task{Function: sampleSumTask, Args: Args{6, 7}}
	id, ack := q.AddTask(task)
	<-ack

	if val, err := store.Get(id); err != nil {
		t.Errorf("Job %s not in state", id)
	} else {
		if val.Status != JobCompleted {
			t.Error("Incorrect status: Expected: Completed, Got: ", val.Status)
		}

		expected := 13
		if val.Value.(int) != expected {
			t.Errorf("Incorrect value: Expected: %d, Got: %d", expected, val.Value)
		}
	}
}

// Checks if handleJob recovers from panics
// In this test, we intentionaly divide an integer be zero to create a panic
// We check if the situation is handled properly
func TestWorkerhandleJobError(t *testing.T) {
	sampleFailingTask := func(args Args) Outcome {
		res := args[0].(int) / args[1].(int)
		return Outcome{Value: res}
	}

	task := Task{Function: sampleFailingTask, Args: Args{10, 0}}

	id, ack := q.AddTask(task)
	<-ack

	if val, err := store.Get(id); err != nil {
		t.Errorf("Job %s not in state", id)
	} else {
		if val.Status != JobError {
			t.Errorf("Incorrect job status: Expected: %s Got: %s", JobError, val.Status)
		}
	}
}
