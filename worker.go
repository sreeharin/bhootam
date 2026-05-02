package bhootam

// handleJob is the job runner
// it's called by the worker goroutine
func handleJob(queue *Queue, store *Store, job Job) {
	defer job.ctxCancel()

	// Acknowledge a worker has picked up the job
	job.ack <- struct{}{}
	store.Set(job.id, Result{Status: JobRunning})

	var (
		outcome Outcome
		status  JobState
	)

	outcomeCh := make(chan Outcome, 1)

	taskComplete := make(chan struct{})
	taskError := make(chan struct{})

	// Run the task in a goroutine
	go func() {
		// handle panic
		defer func() {
			if r := recover(); r != nil {
				taskError <- struct{}{}
			}
		}()

		outcomeCh <- job.task.Run()
		taskComplete <- struct{}{}
	}()

	// Helps handle the timeout funtionality
loop:
	for {
		select {
		case <-job.ctx.Done():
			// If the timeout is reached
			status = JobTimeOut
			outcome = <-outcomeCh
			break loop
		case <-taskComplete:
			// Everything went smoothly
			status = JobCompleted
			outcome = <-outcomeCh
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
		newJob := job
		newJob.task.DecrementRetry()

		queue.jobs <- newJob

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
