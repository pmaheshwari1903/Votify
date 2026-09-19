package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/votify/pkg/logging"
	"github.com/votify/realtime-service/internal/pubsub"
)

type Client struct {
	PollID string
	Conn   *websocket.Conn
	Send   chan pubsub.RealtimeUpdateMessage
}

type ConnectionManager struct {
	mu         sync.RWMutex
	clients    map[string]map[*Client]bool // pollId -> map[Client]bool
	redisStore *pubsub.RedisRealtimeStore
	logger     *logging.Logger
}

func NewConnectionManager(redisStore *pubsub.RedisRealtimeStore, logger *logging.Logger) *ConnectionManager {
	return &ConnectionManager{
		clients:    make(map[string]map[*Client]bool),
		redisStore: redisStore,
		logger:     logger,
	}
}

func (m *ConnectionManager) Register(client *Client) {
	m.mu.Lock()
	if m.clients[client.PollID] == nil {
		m.clients[client.PollID] = make(map[*Client]bool)
	}
	m.clients[client.PollID][client] = true
	m.mu.Unlock()

	m.logger.Info("WebSocket client connected", "poll_id", client.PollID)

	// Subscribe this client to Redis updates for their pollID
	unsub := m.redisStore.Subscribe(client.PollID, func(msg pubsub.RealtimeUpdateMessage) {
		select {
		case client.Send <- msg:
		default:
			m.logger.Warn("Client send buffer full, dropping frame", "poll_id", client.PollID)
		}
	})

	// Start writer loop
	go func() {
		defer func() {
			unsub()
			m.Unregister(client)
			client.Conn.Close()
		}()

		for msg := range client.Send {
			if err := client.Conn.WriteJSON(msg); err != nil {
				m.logger.Warn("Failed to send WebSocket message", "error", err)
				break
			}
		}
	}()
}

func (m *ConnectionManager) Unregister(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clientsMap, exists := m.clients[client.PollID]; exists {
		if _, ok := clientsMap[client]; ok {
			delete(clientsMap, client)
			close(client.Send)
			if len(clientsMap) == 0 {
				delete(m.clients, client.PollID)
			}
			m.logger.Info("WebSocket client disconnected", "poll_id", client.PollID)
		}
	}
}
