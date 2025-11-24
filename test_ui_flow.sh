#!/bin/bash

echo "=== Ephemeral Chat UI Test Flow ==="
echo

echo "1. Making Bob online via backend..."
curl -s -X POST "http://localhost:3000/presence/ping" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"Bob"}' | jq .
echo

echo "2. Now Alice (from UI) can send message to Bob..."
curl -s -X POST "http://localhost:3000/message/send" \
  -H "Content-Type: application/json" \
  -d '{"from":"Alice","to":"Bob","content":"Hello from the UI test!"}' | jq .
echo

echo "3. Bob can receive the message..."
curl -s -X GET "http://localhost:3000/message/receive?user_id=Bob" | jq .
echo

echo "=== UI Flow Instructions ==="
echo "Now in the UI:"
echo "1. Join as 'Alice' in one browser window"
echo "2. Join as 'Bob' in another browser window"
echo "3. Wait a few seconds for presence pings to register"
echo "4. Send messages between Alice and Bob"
echo
echo "Note: Both users need to be actively using the UI for presence to work!"