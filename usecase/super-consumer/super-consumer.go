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

func (c *consumer) Register(ctx context.Context, responses <-chan events.EventResponse) {
	log.Printf("super consumer start reading")

	var size uint

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 30):
				log.Println("total consumed size:", human_readable.ByteCountSI(int64(size)))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-responses:
			size += uint(len(r.Value))
		}
	}
}
