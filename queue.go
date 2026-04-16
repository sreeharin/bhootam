package bhootam

import "github.com/google/uuid"

type Job struct {
	id   string
	task Func
}

type Queue struct {
	jobs chan Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan Job)}
}

func (q *Queue) AddTask(task Func) string {
	id := uuid.NewString()

	q.jobs <- Job{id, task}

	return id
}
