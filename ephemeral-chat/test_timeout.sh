#!/bin/bash

BASE_URL="http://localhost:3000"

echo "=== Testing Presence Timeout ==="
echo

# Test 1: Mark user as online
echo "1. Mark Alice as online:"
curl -s -X POST "$BASE_URL/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice"}' | jq .
echo

# Test 2: Send message immediately (should work)
echo "2. Send message to Alice (should work):"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "bob", "to": "alice", "content": "Hello Alice!"}' | jq .
echo

# Test 3: Wait for timeout (31 seconds > 30 second timeout)
echo "3. Waiting 31 seconds for Alice to go offline..."
sleep 31

# Test 4: Try to send message after timeout (should fail)
echo "4. Try to send message after timeout (should fail):"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "bob", "to": "alice", "content": "Are you still there?"}' | jq .
echo

echo "=== Timeout test complete! ==="