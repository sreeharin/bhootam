package bhootam

import (
	"testing"
)

func sampleTestTask(args Args) Outcome {
	return Outcome{}
}

func TestAddTask(t *testing.T) {
	q := NewQueue()
	// ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	// defer cancel()

	task := NewTask(sampleTestTask)

	go func() {
		q.AddTask(task)
	}()

	job := <-q.jobs
	if job.id == "" {
		t.Errorf("Non empty job id expected")
	}
}
