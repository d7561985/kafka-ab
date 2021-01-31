package super_consumer

import (
	"context"
	"kafka-bench/events"
	"log"
	"time"
)

type worker struct {
	idx   int
	size  *uint
	event <-chan events.EventResponse
	debug bool
}

func (c *worker) work(ctx context.Context) {
	ok := false

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
			if !ok && c.debug {
				log.Printf("consumer[%d] no work", c.idx)
			}
			ok = false
		case r := <-c.event:
			*c.size += uint(len(r.Value))
			ok = true
		}
	}
}
