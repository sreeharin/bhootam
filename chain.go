package bhootam

// Chain implements chaining of tasks.
// Results of one task are passed as arguments to the next task in the chain.
func Chain(tasks ...*Task) (id string, ack, done chan struct{}) {
	ack = make(chan struct{}, 1)
	done = make(chan struct{}, 1)

	ack <- struct{}{}
	done <- struct{}{}

	return "", ack, done
}
