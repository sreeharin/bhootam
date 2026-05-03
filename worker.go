package bhootam

import (
	"context"
	"fmt"
)

// handleJob is the job runner
// it's called by the worker goroutine
func handleJob(queue *Queue, store *Store, job *Job) {
	if !store.Acquire(job.id) {
		return
	}

	defer store.Release(job.id)

	fmt.Println("attempt ", job.task.Retry.Load())

	defer job.ctxCancel()

	// Acknowledge a worker has picked up the job
	job.ack <- struct{}{}
	store.Set(job.id, Result{Status: JobRunning})

	var (
		outcome Outcome
		status  JobState
	)

	outcomeCh := make(chan Outcome, 1)
	taskError := make(chan struct{}, 1)

	// Run the task in a goroutine
	go func() {
		// handle panic
		defer func() {
			if r := recover(); r != nil {
				taskError <- struct{}{}
			}
		}()

		outcomeCh <- job.task.Run()
	}()

	// Helps handle the timeout funtionality
loop:
	for {
		select {
		case <-job.ctx.Done():
			// If the timeout is reached
			status = JobTimeOut
			break loop
		case outcome = <-outcomeCh:
			// Everything went smoothly
			status = JobCompleted
			break loop
		case <-taskError:
			// Signalled from defer func
			status = JobError
			break loop
		}
	}

	// If an error was returned by the user
	// change the JobState to reflect it
	if outcome.Err != nil {
		status = JobError
	}

	// Handle task retry
	if status == JobError && job.task.Retry.Load() > 0 {
		job.task.DecrementRetry()

		ctx, cancel := context.WithTimeout(context.TODO(), job.task.Timeout)

		newJob := Job{
			ctx:       ctx,
			ctxCancel: cancel,
			id:        job.id,
			task:      job.task,
			done:      job.done,
		}

		queue.jobs <- &newJob

		queue.count.Add(1)
		return
	}

	store.Set(job.id, Result{Outcome: outcome, Status: status})
	job.done <- struct{}{}
}

func worker(queue *Queue, store *Store) {
	for job := range queue.jobs {
		queue.count.Add(-1)

		handleJob(queue, store, job)
	}
}

func StartWorker(queue *Queue, store *Store) {
	const numWorkers = 10

	for range numWorkers {
		go worker(queue, store)
	}
}
