package super_consumer

import (
	"context"
	"kafka-bench/events"
	human_readable "kafka-bench/internal/human-readable"
	"kafka-bench/usecase"
	"log"
	"time"
)

type consumer struct {
}

func New() usecase.SuperConsumer {
	return &consumer{}
}

func (c *consumer) Register(ctx context.Context, responses []<-chan events.EventResponse) {
	log.Printf("super consumer start reading with %d threads", len(responses))

	size := make([]uint, len(responses))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 30):
				var v uint
				for _, u := range size {
					v += u
				}

				log.Println("total consumed size:", human_readable.ByteCountSI(int64(v)))
			}
		}
	}()

	for i, r := range responses {
		w := worker{
			idx:   i,
			size:  &size[i],
			event: r,
			debug: true,
		}

		go w.work(ctx)
	}
}
