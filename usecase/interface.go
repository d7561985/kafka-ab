package usecase

import (
	"context"
	"kafka-bench/events"
)

type SuperConsumer interface {
	Register(context.Context, []<-chan events.EventResponse)
}
