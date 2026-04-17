package bhootam

import "testing"

func sampleTestTask(args Args) Outcome {
	return Outcome{}
}

func TestAddTask(t *testing.T) {
	q := NewQueue()

	go func() { q.AddTask(Task{Function: sampleTestTask, Args: Args{}}) }()

	job := <-q.jobs
	if job.id == "" {
		t.Errorf("Non empty job id expected")
	}
}
