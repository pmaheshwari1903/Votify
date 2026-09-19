package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/votify/pkg/events"
	"github.com/votify/pkg/logging"
)

// EventBus is an in-memory event stream bus used for local realtime propagation when Kafka brokers are offline.
type EventBus struct {
	mu        sync.RWMutex
	listeners []func(event events.Event)
}

var DefaultEventBus = &EventBus{}

func (b *EventBus) Subscribe(fn func(event events.Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, fn)
}

func (b *EventBus) Publish(event events.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, fn := range b.listeners {
		go fn(event)
	}
}

// KafkaProducer manages domain event publishing for Vote Service.
type KafkaProducer struct {
	brokers []string
	logger  *logging.Logger
}

func NewKafkaProducer(brokers []string, logger *logging.Logger) *KafkaProducer {
	return &KafkaProducer{
		brokers: brokers,
		logger:  logger,
	}
}

func (p *KafkaProducer) PublishVoteCreated(ctx context.Context, payload events.VoteCreatedPayload) error {
	eventID := fmt.Sprintf("evt_%d", time.Now().UnixNano())
	event := events.NewEvent(
		eventID,
		events.EventTypeVoteCreated,
		"vote-service",
		payload.PollID,
		payload,
	)

	p.logger.Info("Publishing VoteCreated event",
		"event_id", event.EventID,
		"poll_id", payload.PollID,
		"option_id", payload.OptionID,
	)

	// Publish to internal EventBus for instant local realtime fan-out
	DefaultEventBus.Publish(event)

	return nil
}
