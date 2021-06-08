package super_producer

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/d7561985/kafka-ab/adapter"
	human_readable "github.com/d7561985/kafka-ab/internal/human-readable"

	"github.com/icrowley/fake"
)

type producer struct {
	Config
	emitter          adapter.Producer
	messagesNumber   int64
	completeMessages int64
	done             chan struct{}
}

func New(e adapter.Producer, c Config) *producer {
	return &producer{emitter: e, Config: c, done: make(chan struct{})}
}

func (p *producer) Done() <-chan struct{} {
	return p.done
}

// Run blocking function
func (p *producer) Run(ctx context.Context) {
	if p.WindowSize <= 0 || p.Concurrency <= 0 {
		log.Fatalf("window size %d or concurency %d should be greate then or equal 1", p.WindowSize, p.Concurrency)
	}

	if p.Requests > 0 && p.Requests < p.Concurrency {
		log.Printf("decrease concurency %d because of less request target %d", p.Concurrency, p.Requests)
		p.Concurrency = p.Requests
	}

	var wg sync.WaitGroup
	wg.Add(p.Concurrency)

	for i := 0; i < p.Concurrency; i++ {
		go func() {
			p.work(ctx)
			wg.Done()
		}()
	}

	go p.timerInfo(ctx)

	target := "unlimited messages"
	if p.Requests > 0 {
		target = fmt.Sprintf("with target %d message(s)", p.Requests)
	}

	log.Printf("producer worker starts %d concurency with window size %s %s",
		p.Concurrency, human_readable.ByteCountSI(int64(p.WindowSize)), target)

	wg.Wait()

	log.Printf("successfully produced %d messages", p.messagesNumber)
	p.done <- struct{}{}
}

func (p *producer) timerInfo(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			if p.Verbosity < 1 {
				break
			}

			var total = p.completeMessages * int64(p.WindowSize)
			log.Println("msg:", p.completeMessages, "total size:", human_readable.ByteCountSI(total))
		}
	}
}

func (p *producer) work(ctx context.Context) {
	buf := make([]byte, p.WindowSize)
	n, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("init rand error: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// meat target should silentrly leave
			if p.Requests > 0 && atomic.LoadInt64(&p.messagesNumber) >= int64(p.Requests) {
				return
			}

			atomic.AddInt64(&p.messagesNumber, 1)

			if p.do(buf[:n]) != nil {
				atomic.AddInt64(&p.messagesNumber, -1)
				continue
			}

			atomic.AddInt64(&p.completeMessages, 1)
		}
	}
}

func (p *producer) do(v []byte) error {
	// rand key send to rand partition ;)
	key := fake.UserName()

	if err := p.emitter.Emit(p.Topic, []byte(key), v); err != nil {
		log.Fatalf("emit topic %q key %q error: %s", p.Topic, key, err)
		return err
	}

	return nil
}
