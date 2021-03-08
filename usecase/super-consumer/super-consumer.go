package super_consumer

import (
	"context"
	"kafka-bench/events"
	human_readable "kafka-bench/internal/human-readable"
	"kafka-bench/usecase"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type consumer struct {
	cfg  Config
	wg   sync.WaitGroup
	done chan struct{}

	size  uint64
	count uint64

	cancel context.CancelFunc
}

func New(c Config) usecase.SuperConsumer {
	return &consumer{cfg: c, done: make(chan struct{})}
}

func (c *consumer) Start(root context.Context, responses []<-chan events.EventResponse) <-chan struct{} {
	// we can use more then one register call, just add more threads to waight
	c.wg.Add(len(responses))

	log.Printf("super consumer start reading with %d threads", len(responses))

	ctx, cancel := context.WithCancel(root)
	c.cancel = cancel

	for i, r := range responses {
		w := newWorker(i, r)
		go w.work(ctx)
		go c.readWorker(w.size)
	}

	go func() {
		c.wg.Wait()
		c.done <- struct{}{}
	}()

	go c.printDebug(ctx)

	return c.done
}

// readWorker read all size data from worker via channel and if channel is closed this mean that worker finished to consume
func (c *consumer) readWorker(ch <-chan uint) {
	for u := range ch {
		atomic.AddUint64(&c.size, uint64(u))

		if v := atomic.AddUint64(&c.count, 1); c.cfg.Requests > 0 && v >= c.cfg.Requests {
			// cancel will close all workers channels and we weill fallthrough
			c.cancel()
		}
	}

	c.wg.Done()

	log.Printf("consumer: total packets: %d size of: %s", c.count, human_readable.ByteCountSI(int64(c.size)))
}

func (c *consumer) printDebug(ctx context.Context) {
	if c.cfg.Verbosity < 1 {
		return
	}

	for {
		select {
		case <-ctx.Done():
			if c.count == 0 {
				log.Printf("consumer doesn't read any message")
			}
			return
		case <-time.After(time.Second):
			log.Printf("consumer: total packets: %d size of: %s", c.count, human_readable.ByteCountSI(int64(c.size)))
		}
	}
}
