#!/bin/bash

BASE_URL="http://localhost:3000"

echo "=== Testing Ephemeral Chat API ==="
echo

# Test 1: Health check
echo "1. Health check:"
curl -s -X GET "$BASE_URL/health" | jq .
echo

# Test 2: Try to send message to offline user
echo "2. Try sending message to offline user (should fail):"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "bob", "content": "Hello Bob!"}' | jq .
echo

# Test 3: Mark Bob as online
echo "3. Mark Bob as online (presence ping):"
curl -s -X POST "$BASE_URL/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "bob"}' | jq .
echo

# Test 4: Send message to online user (should succeed)
echo "4. Send message to online Bob (should succeed):"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "bob", "content": "Hello Bob!"}' | jq .
echo

# Test 5: Bob receives the message
echo "5. Bob receives message:"
curl -s -X GET "$BASE_URL/message/receive?user_id=bob" | jq .
echo

# Test 6: Try to receive message again (should be empty - ephemeral!)
echo "6. Try to receive message again (should be empty - ephemeral!):"
curl -s -X GET "$BASE_URL/message/receive?user_id=bob" | jq .
echo

# Test 7: Send another message
echo "7. Send another message:"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "bob", "content": "How are you doing?"}' | jq .
echo

# Test 8: Wait a moment and receive
echo "8. Bob receives second message:"
curl -s -X GET "$BASE_URL/message/receive?user_id=bob" | jq .
echo

echo "=== Testing complete! ==="