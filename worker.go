package bhootam

import (
	"context"
	"time"

	"golang.org/x/exp/rand"
)

const (
	Delay    = 500 * time.Millisecond
	MaxDelay = 1 * time.Minute
)

// retryBackoff implements exponential backoff strategy
func retryBackoff(queue *Queue, store *Store, job *Job, attempt int) {
	base := Delay * time.Duration(1<<attempt)

	if base > MaxDelay {
		base = MaxDelay
	}

	delay := time.Duration(rand.Int63n(int64(base)))

	time.Sleep(delay)
	queue.Enqueue(job)
}

// handleJob is the job runner
// it's called by the worker goroutine
func handleJob(id int, queue *Queue, store *Store, job *Job) {
	defer job.ctxCancel()

	// acknowledge a worker has picked up the job
	// for retries it doesn't create an ack
	// might change in the future
	if !job.retry {
		job.ack <- struct{}{}
		close(job.ack)
	}
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

	// Helps handle the status of the job
	// handles various scenarios like timeout, completion, and error
loop:
	for {
		select {
		case <-job.ctx.Done():
			// If the timeout is reached
			status = JobTimeOut
			break loop
		case outcome = <-outcomeCh:
			if outcome.Err != nil {
				status = JobError
			} else {
				// Everything went smoothly
				status = JobCompleted
			}
			break loop
		case <-taskError:
			// Signalled from defer func
			status = JobError
			break loop
		}
	}

	// Handle task retry
	if status == JobError && job.task.retry.Load() > 0 {
		job.task.DecrementRetry()

		ctx, cancel := context.WithTimeout(context.TODO(), job.task.timeout)
		retryJob := NewJob(ctx, cancel, job.id, job.task, withDone(job.done), withJobRetry())

		store.Set(job.id, Result{Status: JobRetry})
		attempt := job.task.maxRetry - int(job.task.retry.Load()+1)
		go retryBackoff(queue, store, retryJob, attempt)

		return
	}

	store.Set(job.id, Result{Outcome: outcome, Status: status})
	job.done <- struct{}{}
}

func worker(id int, queue *Queue, store *Store) {
	for job := range queue.jobs {
		queue.count.Add(-1)

		handleJob(id, queue, store, job)
	}
}

func StartWorker(queue *Queue, store *Store) {
	const numWorkers = 10

	for id := range numWorkers {
		go worker(id, queue, store)
	}
}
