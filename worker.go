package bhootam

// handleJob is the job runner
// it's called by the worker goroutine
func handleJob(store *Store, job Job) {
	defer job.ctxCancel()

	// Acknowledge a worker has picked up the job
	job.ack <- struct{}{}
	store.Set(job.id, Result{Status: JobRunning})

	var (
		outcome Outcome
		status  JobState
	)

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

		outcome = job.task.Run()

		taskComplete <- struct{}{}
	}()

	// Helps handle the timeout funtionality
loop:
	for {
		select {
		case <-job.ctx.Done():
			// If the timeout is reached
			status = JobTimeOut
			break loop
		case <-taskComplete:
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

	store.Set(job.id, Result{Outcome: outcome, Status: status})
	job.done <- struct{}{}
}

func worker(queue *Queue, store *Store) {
	for job := range queue.jobs {
		handleJob(store, job)
	}
}

func StartWorker(queue *Queue, store *Store) {
	const numWorkers = 10

	for range numWorkers {
		go worker(queue, store)
	}
}
