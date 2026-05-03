package bhootam

import (
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

var (
	q     *Queue
	store *Store
)

func TestMain(m *testing.M) {
	q = NewQueue()
	store = NewStore()

	StartWorker(q, store)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func sampleDivideTask(args Args) Outcome {
	res := args[0].(int) / args[1].(int)
	return Outcome{Value: res}
}

func sampleSlowTask(args Args) Outcome {
	for range 3 {
		time.Sleep(1 * time.Second)
	}
	return Outcome{}
}

func TestHandleJob(t *testing.T) {
	tests := []struct {
		name           string
		function       Func
		args           Args
		timeout        time.Duration
		retries        int
		expectedStatus JobState
		expectedValue  any
	}{
		{name: "add job to task and get result", function: sampleSumTask, args: Args{6, 7}, expectedStatus: JobCompleted, expectedValue: 13},
		{name: "check if worker recovers from panic", function: sampleDivideTask, args: Args{10, 0}, expectedStatus: JobError},
		{name: "timeout is respected", function: sampleSlowTask, timeout: 10 * time.Millisecond, expectedStatus: JobTimeOut},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(tt.function, withArgs(tt.args), withTimeout(tt.timeout))

			id, _, done := q.CreateJob(task)
			<-done

			if res, err := store.Get(id); err != nil {
				t.Errorf("Job id: %s not in Store", id)
			} else {
				if tt.expectedStatus != res.Status {
					t.Errorf("Wrong status. Expected: %s, Got: %s", tt.expectedStatus, res.Status)
				}

				if tt.expectedValue != res.Value {
					t.Errorf("Wrong value. Expected: %v Got: %v", tt.expectedValue, res.Value)
				}
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	var retries int32
	retries = 4
	var attempts atomic.Int32

	sampleRetryTask := func(args Args) Outcome {
		attempts.Add(1)
		return Outcome{Err: errors.New("Unexpected error")}
	}

	task := NewTask(sampleRetryTask, withRetry(retries))
	_, ack, done := q.CreateJob(task)
	<-ack

	// Job was taken from the queue by a worker
	if q.count.Load() != 0 {
		t.Errorf("Unexpected job count. Expected: 0, Got: %d", q.count.Load())
	}

	select {
	case <-time.After(5000 * time.Millisecond):
		t.Errorf("Timeout reached. Task wasn't enqueued.")
	case <-done:
		if attempts.Load() != int32(retries+1) {
			t.Errorf("Attempt mismatch. Expected: %d, Got: %d", retries+1, attempts.Load())
		}
	}
}
