package bhootam

import (
	"testing"
)

func sampleTestTask(args Args) Outcome {
	return Outcome{}
}

func TestAddTask(t *testing.T) {
	q := NewQueue()

	task := NewTask(sampleTestTask)

	go func() {
		q.AddTask(task)
	}()

	job := <-q.jobs
	if job.id == "" {
		t.Errorf("Non empty job id expected")
	}
}
