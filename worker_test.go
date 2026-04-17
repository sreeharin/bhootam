package bhootam

import "testing"

func TestStartWorker(t *testing.T) {
	q := NewQueue()
	store := NewStore()

	StartWorker(q, store)

	task := Task{Function: sampleSumTask, Args: Args{6, 7}}
	id := q.AddTask(task)

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
