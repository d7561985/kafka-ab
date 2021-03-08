package usecase

import (
	"context"
	"kafka-bench/events"
)

type SuperConsumer interface {
	Start(context.Context, []<-chan events.EventResponse) <-chan struct{}
}
