package super_consumer

import (
	"context"

	"github.com/d7561985/kafka-ab/events"
)

type worker struct {
	idx   int
	event <-chan events.EventResponse

	// i use non-buffered channel because i need to be aware that we read exact num what we expect
	// it's possible to read more that expected also in some cases, but now i think it's applicable
	size chan uint
}

func newWorker(i int, e <-chan events.EventResponse) *worker {
	return &worker{
		idx:   i,
		event: e,
		size:  make(chan uint),
	}
}

func (c *worker) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(c.size)
			return
		case r := <-c.event:
			c.size <- uint(len(r.Value))
		}
	}
}
