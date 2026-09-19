package broker

import (
	"context"
	"testing"
	"time"

	"github.com/votify/pkg/events"
	"github.com/votify/pkg/logging"
	pollDto "github.com/votify/poll-service/internal/dto"
	pollRepo "github.com/votify/poll-service/internal/repository"
	pollRoutes "github.com/votify/poll-service/internal/routes"
	pollSvc "github.com/votify/poll-service/internal/service"
	"github.com/votify/realtime-service/internal/pubsub"
)

func TestKafkaConsumer_EventDeduplicationAndRealtimeFanout(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewLogger("test-realtime", "development", "debug")

	// 1. Create a test poll in shared repository
	sharedPollRepo := pollRoutes.GetSharedPollRepository()
	pSvc := pollSvc.NewPollService(sharedPollRepo)

	poll, err := pSvc.CreatePoll(ctx, "owner_123", pollDto.CreatePollRequest{
		Question: "Which backend framework is fastest?",
		Options:  []string{"Gin", "Fiber", "Echo"},
	})
	if err != nil {
		t.Fatalf("Failed to create test poll: %v", err)
	}

	redisStore := pubsub.NewRedisRealtimeStore("localhost:6379", "")
	consumer := NewKafkaConsumer("votify-realtime-test", redisStore, logger)

	// Subscribe to Redis PubSub updates for this poll
	receivedUpdates := make(chan pubsub.RealtimeUpdateMessage, 10)
	unsub := redisStore.Subscribe(poll.ID, func(msg pubsub.RealtimeUpdateMessage) {
		receivedUpdates <- msg
	})
	defer unsub()

	// 2. Process first VoteCreated event
	evt1 := events.Event{
		EventID:         "evt_test_001",
		EventType:       events.EventTypeVoteCreated,
		EventVersion:    "1.0",
		ProducerService: "vote-service",
		OccurredAt:      time.Now().UTC().Format(time.RFC3339),
		PartitionKey:    poll.ID,
		Payload: events.VoteCreatedPayload{
			PollID:   poll.ID,
			OptionID: poll.Options[0].ID,
		},
	}

	// Increment vote in poll repo first
	_ = sharedPollRepo.IncrementVote(ctx, poll.ID, poll.Options[0].ID)

	consumer.ProcessEvent(ctx, evt1)

	select {
	case msg := <-receivedUpdates:
		if msg.PollID != poll.ID {
			t.Errorf("Expected pollID %s, got %s", poll.ID, msg.PollID)
		}
		if msg.TotalVotes != 1 {
			t.Errorf("Expected 1 total vote, got %d", msg.TotalVotes)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("Timed out waiting for realtime update")
	}

	// 3. Process Duplicate Event (Same EventID)
	consumer.ProcessEvent(ctx, evt1)

	select {
	case <-receivedUpdates:
		t.Fatalf("Expected duplicate event to be ignored, but received update!")
	case <-time.After(200 * time.Millisecond):
		// Success — duplicate was properly ignored
	}
}
