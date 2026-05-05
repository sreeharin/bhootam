# Bhootam

Bhootam is an experimental, in-memory task queue library written in Go. It's designed as a toy project for learning and experimentation, not intended for production use. It provides basic task queuing functionality with support for timeouts, retries, and panic recovery.

## Features

- **In-Memory Queue**: Simple channel-based job queue for asynchronous task execution.
- **Timeout Support**: Configurable timeouts for tasks to prevent hanging operations.
- **Retry Mechanism**: Automatic retry for failed jobs with configurable retry counts.
- **Panic Recovery**: Workers recover from panics in task functions.

## Basic usage

### Installation
```
go get github.com/sreeharin/bhootam@latest
```

### Task
The task should follow a specific function signature
```
func (args bm.Args) bm.Outcome
```

```
func sum(args bm.Args) bm.Outcome {
    res := args[0].(int) + args[1].(int)
    return bm.Outcome{Value: res}
}
```
Note: Since args accept arbitary number of arguments of all type, arguments should be type casted.

### Queue, Store, and Workers
Since bhootam is an in-memory task queue, make sure to initialize the queue, and the store.
Also we've to start the workers.

```
queue := bm.NewQueue()
store := bm.NewStore()

bm.StartWorker(queue, store)
```

### Example
```
import (
	"fmt"

	bm "github.com/sreeharin/bhootam"
)

func main() {
	sum := func(args bm.Args) bm.Outcome {
		res := args[0].(int) + args[1].(int)
		return bm.Outcome{Value: res}
	}

	queue := bm.NewQueue()
	store := bm.NewStore()

	bm.StartWorker(queue, store)

	task := bm.NewTask(sum, bm.WithArgs(bm.Args{5, 5}))
	id, ack, _ := queue.CreateJob(task)

	<-ack

	res, _ := store.Get(id)

	fmt.Println(res.Outcome.Value) // 10
}
```

## License

See LICENSE file for details.
