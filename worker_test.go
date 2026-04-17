package bhootam

import "testing"

func TestStartWorker(t *testing.T) {
	q := NewQueue()
	store := NewStore()

	StartWorker()

	task := Task{Function: sampleSumTask, Args: Args{6, 7}}
	id := q.AddTask(task)

	if val, err := store.Get(id); err != nil {
		t.Errorf("Job %s not in state", id)
	} else {
		if val.Status != JobCompleted {
			t.Error("Incorrect status: Expected: Completed, Got: ", val.Status)
		}

		if val.Outcome.Value.(int) != 13 {
			t.Error("Incorrect value: Expected: 13, Got: ", val.Value)
		}

	}

}
