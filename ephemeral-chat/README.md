# Ephemeral Chat API

A minimal-traffic, REST-only, ephemeral chat system where messages exist only briefly after being fetched, with no message history stored on the server.

## Features

- **Ephemeral Messages**: Messages are deleted immediately after being delivered
- **Presence-Based**: Users must be online to receive messages
- **Pull-Only Delivery**: No push notifications, recipients must poll for messages
- **No Persistence**: Everything lives in memory, no databases
- **Auto-Cleanup**: Expired messages and offline users are automatically cleaned up

## Configuration

```go
const (
    PRESENCE_TIMEOUT = 30 * time.Second // User goes offline after 30s without ping
    MESSAGE_TTL      = 60 * time.Second // Messages expire after 60s if not retrieved
)
```

## API Endpoints

### Health Check
```
GET /health
```
Returns server status.

### Presence Management

#### Ping Presence
```
POST /presence/ping
Content-Type: application/json

{
    "user_id": "string"
}
```
Marks user as online. Users must ping periodically to stay online.

**Response:**
```json
{
    "status": "ok",
    "user_id": "string",
    "timestamp": "2025-11-24T09:37:47.868848+01:00"
}
```

### Messaging

#### Send Message
```
POST /message/send
Content-Type: application/json

{
    "from": "string",
    "to": "string",
    "content": "string"
}
```

**Success Response:**
```json
{
    "status": "sent",
    "message_id": "uuid",
    "timestamp": "2025-11-24T09:37:47.880008+01:00"
}
```

**Error Response (User Offline):**
```json
{
    "error": "user_offline"
}
```

#### Receive Messages
```
GET /message/receive?user_id=string
```

**Response (Message Available):**
```json
{
    "message": {
        "id": "uuid",
        "from": "string",
        "to": "string",
        "content": "string",
        "timestamp": "2025-11-24T09:37:47.880008+01:00",
        "expires_at": "2025-11-24T09:38:47.880008+01:00"
    }
}
```

**Response (No Messages):**
```json
{
    "messages": []
}
```

## How It Works

1. **Presence Tracking**: Users must call `/presence/ping` periodically to stay online
2. **Message Sending**: Sender calls `/message/send` - fails if recipient is offline
3. **Message Receiving**: Recipient calls `/message/receive` to get one message
4. **Ephemeral Nature**: Once a message is delivered, it's immediately deleted
5. **Auto-Cleanup**: Background routines clean up expired messages and offline users

## Example Usage

See `test_api.sh` for a complete example flow:

```bash
# 1. Try sending to offline user (fails)
curl -X POST "http://localhost:3000/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "bob", "content": "Hello!"}'

# 2. Mark recipient as online
curl -X POST "http://localhost:3000/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "bob"}'

# 3. Send message (succeeds)
curl -X POST "http://localhost:3000/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "bob", "content": "Hello!"}'

# 4. Receive message
curl -X GET "http://localhost:3000/message/receive?user_id=bob"

# 5. Try to receive again (empty - message is gone!)
curl -X GET "http://localhost:3000/message/receive?user_id=bob"
```

## Running the Server

```bash
go build -o ephemeral-chat-server .
./ephemeral-chat-server
```

Server runs on port 3000 by default.

## Key Characteristics

- **No Message History**: Messages don't persist anywhere
- **Real-Time Requirement**: Both users must be active for communication
- **No Delivery Guarantees**: If recipient doesn't poll, message expires
- **Memory-Only**: Restart loses all data (by design)
- **Single Message Pull**: Each request returns max one message

This creates a truly ephemeral, "in-person style" chat experience where timing matters and nothing is permanently stored.