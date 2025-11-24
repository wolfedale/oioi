package main

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	PRESENCE_TIMEOUT     = 90 * time.Second  // User goes offline after 90s without ping
	MESSAGE_SEND_TIMEOUT = 60 * time.Second  // Allow message sending if user was online within 60s
	MESSAGE_TTL          = 60 * time.Second  // Messages expire after 60s if not retrieved
	OFFLINE_MESSAGE_DROP = 300 * time.Second // Drop all messages for users offline > 5 minutes
)

// User represents a user in the system
type User struct {
	ID       string    `json:"id"`
	LastSeen time.Time `json:"last_seen"`
}

// Message represents an ephemeral message
type Message struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Account represents a user account with credentials
type Account struct {
	UserID       string `json:"user_id"`
	SentenceHash string `json:"-"` // Don't expose in JSON
	CreatedAt    time.Time `json:"created_at"`
}

// ChatServer holds all the in-memory data
type ChatServer struct {
	// Presence tracking
	presence map[string]*User
	presenceMu sync.RWMutex

	// Message inbox for each user
	inboxes map[string][]*Message
	inboxMu sync.RWMutex

	// User accounts (userID -> Account)
	accounts map[string]*Account
	accountsMu sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	server := &ChatServer{
		presence: make(map[string]*User),
		inboxes:  make(map[string][]*Message),
		accounts: make(map[string]*Account),
	}

	// Start cleanup routine
	go server.cleanupRoutine()

	return server
}

// IsUserOnline checks if a user is currently online
func (cs *ChatServer) IsUserOnline(userID string) bool {
	cs.presenceMu.RLock()
	defer cs.presenceMu.RUnlock()

	user, exists := cs.presence[userID]
	if !exists {
		return false
	}

	return time.Since(user.LastSeen) <= PRESENCE_TIMEOUT
}

// IsUserOnlineForMessaging checks if a user is online for message delivery (more lenient)
func (cs *ChatServer) IsUserOnlineForMessaging(userID string) bool {
	cs.presenceMu.RLock()
	defer cs.presenceMu.RUnlock()

	user, exists := cs.presence[userID]
	if !exists {
		return false
	}

	return time.Since(user.LastSeen) <= MESSAGE_SEND_TIMEOUT
}

// UpdatePresence updates user's last seen timestamp
func (cs *ChatServer) UpdatePresence(userID string) {
	cs.presenceMu.Lock()
	defer cs.presenceMu.Unlock()

	cs.presence[userID] = &User{
		ID:       userID,
		LastSeen: time.Now(),
	}
}

// cleanupRoutine runs periodically to clean up expired data
func (cs *ChatServer) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cs.cleanupExpiredMessages()
		cs.cleanupOfflineUsers()
	}
}

// cleanupExpiredMessages removes expired messages from all inboxes
func (cs *ChatServer) cleanupExpiredMessages() {
	cs.inboxMu.Lock()
	cs.presenceMu.RLock()
	defer cs.inboxMu.Unlock()
	defer cs.presenceMu.RUnlock()

	now := time.Now()
	for userID, messages := range cs.inboxes {
		// Check if user has been offline for 5+ minutes
		user, exists := cs.presence[userID]
		isLongOffline := !exists || now.Sub(user.LastSeen) > OFFLINE_MESSAGE_DROP

		if isLongOffline {
			// Drop all messages for users offline > 5 minutes
			cs.inboxes[userID] = make([]*Message, 0)
			continue
		}

		// For online/recently online users, only remove expired messages
		filtered := make([]*Message, 0)
		for _, msg := range messages {
			if now.Before(msg.ExpiresAt) {
				filtered = append(filtered, msg)
			}
		}
		cs.inboxes[userID] = filtered
	}
}

// cleanupOfflineUsers removes users who have been offline too long
func (cs *ChatServer) cleanupOfflineUsers() {
	cs.presenceMu.Lock()
	defer cs.presenceMu.Unlock()

	now := time.Now()
	for userID, user := range cs.presence {
		if now.Sub(user.LastSeen) > PRESENCE_TIMEOUT*2 { // Clean up after 2x timeout
			delete(cs.presence, userID)
		}
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize chat server
	chatServer := NewChatServer()

	// Routes will be added here
	setupRoutes(app, chatServer)

	log.Println("Starting ephemeral chat server on :3000")
	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App, cs *ChatServer) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Account routes
	account := app.Group("/account")
	account.Post("/create", handleAccountCreate(cs))
	account.Post("/login", handleAccountLogin(cs))

	// Presence routes
	presence := app.Group("/presence")
	presence.Post("/ping", handlePresencePing(cs))
	presence.Get("/online", handleGetOnlineUsers(cs))

	// Message routes
	message := app.Group("/message")
	message.Post("/send", handleMessageSend(cs))
	message.Get("/receive", handleMessageReceive(cs))
}