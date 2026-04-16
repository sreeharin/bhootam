package bhootam

import "testing"

func sampleTask(args Args) Value {
	return Value{}
}

func TestAddTask(t *testing.T) {
	q := NewQueue()

	go func() { q.AddTask(sampleTask) }()

	job := <-q.jobs
	if job.id == "" {
		t.Errorf("Non empty job id expected")
	}
}
