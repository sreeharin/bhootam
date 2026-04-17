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
}
