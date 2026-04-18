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
		if val.Outcome.Value.(int) != expected {
			t.Errorf("Incorrect value: Expected: %d, Got: %d", expected, val.Value)
		}
	}
}
