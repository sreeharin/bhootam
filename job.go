package bhootam

type JobState string

const (
	JobCompleted JobState = "completed"
	JobRunning   JobState = "running"
	JobError     JobState = "error"
)

type Job struct {
	id   string
	task Task

	// ack is the acknowledgment provided by the worker
	// to inform that a worker has picken up the task
	ack chan struct{}
}
