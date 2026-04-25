package bhootam

import (
	"context"
	"testing"
	"time"
)

func sampleTestTask(args Args) Outcome {
	return Outcome{}
}

func TestAddTask(t *testing.T) {
	q := NewQueue()
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	go func() { q.AddTask(ctx, Task{Function: sampleTestTask, Args: Args{}}) }()

	job := <-q.jobs
	if job.id == "" {
		t.Errorf("Non empty job id expected")
	}
}
