package broker

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/votify/pkg/events"
	"github.com/votify/pkg/logging"
	"github.com/votify/realtime-service/internal/pubsub"
)

type KafkaConsumer struct {
	groupID       string
	redisStore    *pubsub.RedisRealtimeStore
	logger        *logging.Logger
	mu            sync.Mutex
	processedEvts map[string]bool // Deduplication map for event_id
}

func NewKafkaConsumer(groupID string, redisStore *pubsub.RedisRealtimeStore, logger *logging.Logger) *KafkaConsumer {
	consumer := &KafkaConsumer{
		groupID:       groupID,
		redisStore:    redisStore,
		logger:        logger,
		processedEvts: make(map[string]bool),
	}

	// Subscribe to internal EventBus for immediate local event consumption

	return consumer
}

func (c *KafkaConsumer) ProcessEvent(ctx context.Context, evt events.Event) {
	if evt.EventType != events.EventTypeVoteCreated {
		return
	}

	// 1. Deduplication check on event_id to prevent double processing of Kafka re-deliveries
	c.mu.Lock()
	if c.processedEvts[evt.EventID] {
		c.mu.Unlock()
		c.logger.Warn("Ignoring duplicate VoteCreated event", "event_id", evt.EventID)
		return
	}
	c.processedEvts[evt.EventID] = true
	c.mu.Unlock()

	// 2. Decode payload
	var pollID string
	switch p := evt.Payload.(type) {
	case events.VoteCreatedPayload:
		pollID = p.PollID
	case map[string]interface{}:
		if id, ok := p["poll_id"].(string); ok {
			pollID = id
		}
	default:
		// Attempt JSON unmarshal
		bytes, err := json.Marshal(evt.Payload)
		if err == nil {
			var payload events.VoteCreatedPayload
			if err := json.Unmarshal(bytes, &payload); err == nil {
				pollID = payload.PollID
			}
		}
	}

	if pollID == "" {
		c.logger.Error("Failed to extract poll_id from VoteCreated event", "event_id", evt.EventID)
		return
	}

	c.logger.Info("Realtime Consumer processed VoteCreated event",
		"event_id", evt.EventID,
		"poll_id", pollID,
	)

	// 3. Trigger Redis update & channel publication
	c.redisStore.PublishUpdate(ctx, pollID)
}
