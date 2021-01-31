package super_producer

import (
	"context"
	"crypto/rand"
	"github.com/icrowley/fake"
	"kafka-bench/adapter"
	human_readable "kafka-bench/internal/human-readable"
	"log"
	"time"
)

type producer struct {
	emitter adapter.Producer
	topic   string
	threads int
	msgSize int //

	stat []uint
}

func New(e adapter.Producer, topic string, threads int, msgSize int) *producer {
	return &producer{emitter: e,
		topic:   topic,
		threads: threads, msgSize: msgSize}
}

func (p *producer) Run(ctx context.Context) {
	p.stat = make([]uint, p.threads)

	for i := 0; i < p.threads; i++ {
		go p.work(ctx, &p.stat[i])
	}

	go p.timerInfo(ctx)

	log.Printf("producer worker starts %d threads with msg size %s\n",
		p.threads, human_readable.ByteCountSI(int64(p.msgSize)))
}

func (p *producer) timerInfo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 30):
			var total uint
			for _, v := range p.stat {
				total += v
			}

			log.Println(human_readable.ByteCountSI(int64(total)))
		}
	}
}

func (p *producer) work(ctx context.Context, counter *uint) {
	buf := make([]byte, p.msgSize)
	n, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("init rand error: %s", err)
	}

	size := uint(n)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if p.do(buf[:n]) == nil {
				*counter += size
			}
		}
	}
}

func (p *producer) do(v []byte) error {
	// rand key send to rand partition ;)
	key := fake.UserName()

	if err := p.emitter.Emit(p.topic, []byte(key), v); err != nil {
		log.Printf("[ERROR]: emit topic %q key %q error: %s",
			p.topic, key, err)
		return err
	}

	return nil
}
