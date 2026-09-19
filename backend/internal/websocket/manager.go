package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/votify/backend/internal/realtime"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for WebSocket connections
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	PollID string
	Conn   *websocket.Conn
	Send   chan realtime.RealtimeUpdateMessage
}

type ConnectionManager struct {
	mu              sync.RWMutex
	clients         map[string]map[*Client]bool
	realtimeService *realtime.RealtimeService
}

func NewConnectionManager(realtimeService *realtime.RealtimeService) *ConnectionManager {
	return &ConnectionManager{
		clients:         make(map[string]map[*Client]bool),
		realtimeService: realtimeService,
	}
}

func (m *ConnectionManager) HandleWS(c *gin.Context) {
	pollID := c.Param("pollId")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pollId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		PollID: pollID,
		Conn:   conn,
		Send:   make(chan realtime.RealtimeUpdateMessage, 16),
	}

	m.Register(client)
}

func (m *ConnectionManager) Register(client *Client) {
	m.mu.Lock()
	if m.clients[client.PollID] == nil {
		m.clients[client.PollID] = make(map[*Client]bool)
	}
	m.clients[client.PollID][client] = true
	m.mu.Unlock()

	// Subscribe client to Realtime update notifications
	unsub := m.realtimeService.Subscribe(client.PollID, func(msg realtime.RealtimeUpdateMessage) {
		select {
		case client.Send <- msg:
		default:
			// Buffer full, drop frame
		}
	})

	// Send initial snapshot snapshot frame immediately on connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	snapshot, err := m.realtimeService.BuildSnapshot(ctx, client.PollID)
	cancel()
	if err == nil && snapshot != nil {
		select {
		case client.Send <- *snapshot:
		default:
		}
	}

	// Writer loop
	go func() {
		defer func() {
			unsub()
			m.Unregister(client)
			client.Conn.Close()
		}()

		for msg := range client.Send {
			_ = client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteJSON(msg); err != nil {
				break
			}
		}
	}()

	// Reader loop (keep connection open until disconnect/ping failure)
	go func() {
		defer func() {
			client.Conn.Close()
		}()
		for {
			_, _, err := client.Conn.ReadMessage()
			if err != nil {
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
		}
	}
}
