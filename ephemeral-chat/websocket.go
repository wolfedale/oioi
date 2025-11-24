package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (development only)
		return true
	},
}

// WebSocket client
type WSClient struct {
	UserID     string
	Conn       *websocket.Conn
	Send       chan WSMessage
	ChatServer *ChatServer
}

// WebSocket message types
type WSMessage struct {
	Type      string `json:"type"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Content   string `json:"content,omitempty"`
	MessageID string `json:"message_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Error     string `json:"error,omitempty"`
}

// WebSocket hub manages all clients
type WSHub struct {
	clients    map[string]*WSClient // userID -> client
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan WSMessage
	mu         sync.RWMutex
}

// Create new WebSocket hub
func newWSHub() *WSHub {
	return &WSHub{
		clients:    make(map[string]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan WSMessage),
	}
}

// Run the WebSocket hub
func (h *WSHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("🔌 WebSocket client connected: %s (total: %d)", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("🔌 WebSocket client disconnected: %s (total: %d)", client.UserID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			if targetClient, ok := h.clients[message.To]; ok {
				select {
				case targetClient.Send <- message:
					log.Printf("📨 Message delivered via WebSocket: %s -> %s", message.From, message.To)
				default:
					// Client's send channel is full, remove client
					close(targetClient.Send)
					delete(h.clients, message.To)
					log.Printf("⚠️ Removed unresponsive client: %s", message.To)
				}
			} else {
				log.Printf("❌ Target user not connected via WebSocket: %s", message.To)
			}
			h.mu.RUnlock()
		}
	}
}

// Send message to specific user
func (h *WSHub) sendToUser(userID string, message WSMessage) {
	h.broadcast <- WSMessage{
		Type:      message.Type,
		From:      message.From,
		To:        userID,
		Content:   message.Content,
		MessageID: message.MessageID,
		Timestamp: message.Timestamp,
		ExpiresAt: message.ExpiresAt,
	}
}

// Handle WebSocket connections
func handleWebSocket(hub *WSHub, cs *ChatServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id parameter required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		client := &WSClient{
			UserID:     userID,
			Conn:       conn,
			Send:       make(chan WSMessage, 256),
			ChatServer: cs,
		}

		// Register client
		hub.register <- client

		// Start goroutines
		go client.writePump(hub)
		go client.readPump(hub)
	}
}

// Read messages from WebSocket
func (c *WSClient) readPump(hub *WSHub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	// Set read deadline and pong handler
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle different message types
		switch msg.Type {
		case "send_message":
			c.handleSendMessage(hub, msg)
		case "ping":
			// Update presence
			c.ChatServer.UpdatePresence(c.UserID)
			// Send pong back
			c.Send <- WSMessage{Type: "pong"}
		}
	}
}

// Write messages to WebSocket
func (c *WSClient) writePump(hub *WSHub) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Handle send message request
func (c *WSClient) handleSendMessage(hub *WSHub, msg WSMessage) {
	if msg.To == "" || msg.Content == "" {
		c.Send <- WSMessage{
			Type:  "error",
			Error: "to and content are required",
		}
		return
	}

	// Create message with ID and timestamp
	message := &Message{
		ID:        generateMessageID(),
		From:      c.UserID,
		To:        msg.To,
		Content:   msg.Content,
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(MESSAGE_TTL),
	}

	log.Printf("📤 WebSocket message: %s -> %s: \"%s\"", message.From, message.To, message.Content)

	// Send confirmation to sender
	c.Send <- WSMessage{
		Type:      "message_sent",
		MessageID: message.ID,
		Timestamp: message.Timestamp.Format(time.RFC3339),
	}

	// Send message to recipient via WebSocket
	hub.sendToUser(msg.To, WSMessage{
		Type:      "new_message",
		From:      message.From,
		To:        message.To,
		Content:   message.Content,
		MessageID: message.ID,
		Timestamp: message.Timestamp.Format(time.RFC3339),
		ExpiresAt: message.ExpiresAt.Format(time.RFC3339),
	})
}

// Generate unique message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + generateUserID()
}