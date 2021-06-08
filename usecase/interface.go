package usecase

import (
	"context"

	"github.com/d7561985/kafka-ab/events"
)

type SuperConsumer interface {
	Start(context.Context, []<-chan events.EventResponse) <-chan struct{}
}
