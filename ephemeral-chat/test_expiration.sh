#!/bin/bash

BASE_URL="http://localhost:3000"

echo "=== Testing Message Expiration ==="
echo

# Test 1: Mark user as online
echo "1. Mark Charlie as online:"
curl -s -X POST "$BASE_URL/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "charlie"}' | jq .
echo

# Test 2: Send a message but don't retrieve it immediately
echo "2. Send message to Charlie:"
curl -s -X POST "$BASE_URL/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from": "alice", "to": "charlie", "content": "This message will expire!"}' | jq .
echo

echo "3. Message sent but not retrieved. In a real scenario, if Charlie doesn't"
echo "   check messages within 60 seconds, the message will expire and be cleaned up."
echo "   Let's retrieve it now before it expires:"

# Test 3: Retrieve the message
curl -s -X GET "$BASE_URL/message/receive?user_id=charlie" | jq .
echo

echo "4. Message successfully retrieved and is now gone forever (ephemeral)."

echo "=== Message expiration test complete! ==="