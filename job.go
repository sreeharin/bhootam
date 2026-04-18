package bhootam

type JobState int

const (
	JobCompleted JobState = iota
	JobRunning
	JobError
)

type Job struct {
	id   string
	task Task

	// ack is the acknowledgment provided by the worker
	// to inform that a worker has picken up the task
	ack chan struct{}
}
