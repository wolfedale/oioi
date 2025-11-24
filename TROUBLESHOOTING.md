# Ephemeral Chat - Troubleshooting Guide

## ❌ "Cannot send messages via UI"

### Problem
Messages fail to send through the web interface.

### Root Cause
The ephemeral chat system requires **both users to be actively online** for messaging to work. The frontend shows mock users in the dropdown, but only users who are actively using the interface are truly "online" in the backend.

### Solution Steps

#### 1. Verify Both Servers Are Running
```bash
# Check backend (port 3000)
curl http://localhost:3000/health
# Should return: {"status":"ok"}

# Check frontend (port 3001)
curl http://localhost:3001/
# Should return: HTML page with "Ephemeral Chat"
```

#### 2. Test with Two Browser Windows
1. **Window 1**: Go to http://localhost:3001
   - Join as "Alice"
   - Wait 2-3 seconds for presence ping

2. **Window 2**: Go to http://localhost:3001
   - Join as "Bob"
   - Wait 2-3 seconds for presence ping

3. **Test Messaging**:
   - In Alice's window: Select "Bob" as recipient, send message
   - In Bob's window: Message should appear automatically
   - In Bob's window: Select "Alice" as recipient, reply
   - Messages should flow both ways

#### 3. Common Issues & Fixes

**Issue**: "User is offline" error
- **Cause**: Recipient hasn't joined the chat or their presence expired
- **Fix**: Ensure both users have joined within the last 30 seconds

**Issue**: Messages not appearing
- **Cause**: Message polling not working or user not truly online
- **Fix**: Refresh browser, rejoin chat, wait for presence sync

**Issue**: Can't select recipient
- **Cause**: No other users online
- **Fix**: Open second browser window, join as different user

#### 4. Manual Testing via API

If UI isn't working, test the backend directly:

```bash
# Make Bob online
curl -X POST "http://localhost:3000/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"Bob"}'

# Send message from Alice to Bob
curl -X POST "http://localhost:3000/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from":"Alice","to":"Bob","content":"Test message"}'

# Bob receives message
curl -X GET "http://localhost:3000/message/receive?user_id=Bob"
```

## 🔍 Advanced Debugging

### Check Browser Console
1. Open browser dev tools (F12)
2. Check Console tab for errors
3. Check Network tab for failed API calls

### Monitor Backend Logs
The Go server shows all API requests:
```bash
# Watch for requests in terminal running the server
./ephemeral-chat-server
```

Look for:
- `200` status = success
- `400` status = bad request
- `404` status = user offline

### Presence Timing
- Users ping every 25 seconds
- Users go offline after 30 seconds without ping
- Users are cleaned up after 60 seconds offline

## 🚀 Expected Behavior

### Normal Flow
1. **User joins** → Presence ping starts every 25s
2. **User selects recipient** → Only online users shown
3. **User sends message** → API checks recipient online
4. **Message delivered** → Recipient polls every 2s
5. **Message received** → Appears in UI, deleted from server

### Error Cases
- **Send to offline user** → "user_offline" error
- **Empty message** → Validation error
- **Network issues** → Connection error
- **Server down** → API failure

## 💡 Pro Tips

1. **Keep browser windows active** - inactive tabs may slow polling
2. **Use real names** - easier to test with "Alice" and "Bob"
3. **Watch timing** - 30-second offline timeout is strict
4. **Check CORS** - different ports require CORS (already enabled)
5. **Use dev tools** - Network tab shows all API calls

## 🔧 Quick Fix Commands

```bash
# Restart both servers
cd ephemeral-chat && ./ephemeral-chat-server &
cd ephemeral-chat-frontend && npm run dev &

# Test API health
curl http://localhost:3000/health && curl http://localhost:3001/

# Make test users online
curl -X POST "http://localhost:3000/presence/ping" -H "Content-Type: application/json" -d '{"user_id":"Alice"}'
curl -X POST "http://localhost:3000/presence/ping" -H "Content-Type: application/json" -d '{"user_id":"Bob"}'
```

The system is working correctly - it just requires both users to be actively present! 🎉