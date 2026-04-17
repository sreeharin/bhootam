package bhootam

func worker(queue *Queue, store *Store) {
	for job := range queue.jobs {
		store.Set(job.id, Result{Status: JobRunning})
		outcome := job.task.Run()
		store.Set(job.id, Result{Outcome: outcome, Status: JobCompleted})
	}
}

func StartWorker(queue *Queue, store *Store) {
	const numWorkers = 10

	for range numWorkers {
		go worker(queue, store)
	}

}
