package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PresencePingRequest represents the request for presence ping
type PresencePingRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// MessageSendRequest represents the request for sending a message
type MessageSendRequest struct {
	From    string `json:"from" validate:"required"`
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// MessageReceiveRequest represents the request for receiving messages
type MessageReceiveRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// AccountCreateRequest represents the request for creating an account
type AccountCreateRequest struct {
	Sentence string `json:"sentence" validate:"required"`
}

// AccountLoginRequest represents the request for logging into an account
type AccountLoginRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Sentence string `json:"sentence" validate:"required"`
}

// handlePresencePing handles POST /presence/ping
func handlePresencePing(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PresencePingRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.UserID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user_id is required",
			})
		}

		// Update user presence
		cs.UpdatePresence(req.UserID)

		return c.JSON(fiber.Map{
			"status":    "ok",
			"user_id":   req.UserID,
			"timestamp": time.Now(),
		})
	}
}

// handleMessageSend handles POST /message/send
func handleMessageSend(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req MessageSendRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.From == "" || req.To == "" || req.Content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "from, to, and content are required",
			})
		}

		// Always allow message sending - if recipient is offline, message will be dropped by cleanup routine

		// Create message
		message := &Message{
			ID:        uuid.New().String(),
			From:      req.From,
			To:        req.To,
			Content:   req.Content,
			Timestamp: time.Now(),
			ExpiresAt: time.Now().Add(MESSAGE_TTL),
		}

		// Add ONLY to recipient's inbox (sender will see immediately in UI)
		cs.inboxMu.Lock()
		if cs.inboxes[req.To] == nil {
			cs.inboxes[req.To] = make([]*Message, 0)
		}
		cs.inboxes[req.To] = append(cs.inboxes[req.To], message)
		cs.inboxMu.Unlock()

		return c.JSON(fiber.Map{
			"status":     "sent",
			"message_id": message.ID,
			"timestamp":  message.Timestamp,
		})
	}
}

// handleMessageReceive handles GET /message/receive
func handleMessageReceive(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Query("user_id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user_id query parameter is required",
			})
		}

		// Update presence since user is actively checking messages
		cs.UpdatePresence(userID)

		cs.inboxMu.Lock()
		defer cs.inboxMu.Unlock()

		inbox := cs.inboxes[userID]
		if len(inbox) == 0 {
			return c.JSON(fiber.Map{
				"messages": []Message{},
			})
		}

		// Sort messages by timestamp to ensure correct order
		sort.Slice(inbox, func(i, j int) bool {
			return inbox[i].Timestamp.Before(inbox[j].Timestamp)
		})

		// Get the oldest message (FIFO)
		message := inbox[0]

		// Remove message from inbox (ephemeral - once delivered, it's gone)
		cs.inboxes[userID] = inbox[1:]

		return c.JSON(fiber.Map{
			"message": message,
		})
	}
}

// handleGetOnlineUsers handles GET /presence/online
func handleGetOnlineUsers(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cs.presenceMu.RLock()
		defer cs.presenceMu.RUnlock()

		// Use map to automatically deduplicate users
		onlineUserMap := make(map[string]bool)
		now := time.Now()

		for userID, user := range cs.presence {
			if now.Sub(user.LastSeen) <= PRESENCE_TIMEOUT {
				onlineUserMap[userID] = true
			}
		}

		// Convert map keys to slice
		onlineUsers := make([]string, 0, len(onlineUserMap))
		for userID := range onlineUserMap {
			onlineUsers = append(onlineUsers, userID)
		}

		return c.JSON(fiber.Map{
			"online_users": onlineUsers,
		})
	}
}

// generateUserID creates a unique 8-character ID (not deterministic)
func generateUserID() string {
	// Generate a unique UUID and take first 8 characters
	id := uuid.New().String()
	return id[:8]
}

// generateSentenceHash creates a hash from the sentence for verification
func generateSentenceHash(sentence string) string {
	hash := sha256.Sum256([]byte(sentence))
	hexStr := hex.EncodeToString(hash[:])
	return hexStr[:32] // Use 32 chars for better security
}

// verifySentence checks if a sentence matches the expected hash for a user
func verifySentence(sentence string, expectedHash string) bool {
	// Create SHA256 hash of the sentence
	actualHash := generateSentenceHash(sentence)
	return actualHash == expectedHash
}

// handleAccountCreate handles POST /account/create
func handleAccountCreate(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req AccountCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.Sentence == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "sentence is required",
			})
		}

		// Validate sentence length (minimum reasonable length)
		if len(req.Sentence) < 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "sentence too short (minimum 10 characters)",
			})
		}

		// Generate unique user ID (not deterministic from sentence)
		userID := generateUserID()

		// Create sentence hash for verification
		sentenceHash := generateSentenceHash(req.Sentence)

		// Store account
		cs.accountsMu.Lock()
		cs.accounts[userID] = &Account{
			UserID:       userID,
			SentenceHash: sentenceHash,
			CreatedAt:    time.Now(),
		}
		cs.accountsMu.Unlock()

		return c.JSON(fiber.Map{
			"user_id": userID,
			"message": "Account created successfully",
		})
	}
}

// handleAccountLogin handles POST /account/login
func handleAccountLogin(cs *ChatServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req AccountLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.UserID == "" || req.Sentence == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user_id and sentence are required",
			})
		}

		// Validate sentence length
		if len(req.Sentence) < 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "sentence too short (minimum 10 characters)",
			})
		}

		// Check if account exists
		cs.accountsMu.RLock()
		account, exists := cs.accounts[req.UserID]
		cs.accountsMu.RUnlock()

		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		// Verify sentence matches stored hash
		if !verifySentence(req.Sentence, account.SentenceHash) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		return c.JSON(fiber.Map{
			"user_id": req.UserID,
			"message": "Login successful",
		})
	}
}