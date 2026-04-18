package bhootam

func handleJob(store *Store, job Job) {
	// handle panics
	defer func() {
		if r := recover(); r != nil {
			store.Set(job.id, Result{Status: JobError})
		}
	}()

	// Acknowledge a worker has picked up the job
	job.ack <- struct{}{}

	store.Set(job.id, Result{Status: JobRunning})

	outcome := job.task.Run()

	// If an error was returned by the user
	// change the JobState to reflect it
	var status JobState
	if outcome.Err != nil {
		status = JobError
	} else {
		status = JobCompleted
	}

	store.Set(job.id, Result{Outcome: outcome, Status: status})
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
