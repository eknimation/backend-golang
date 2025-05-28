#!/bin/bash

# JWT Authentication API Test Script
# This script demonstrates the complete authentication flow

# Configuration
BASE_URL="http://localhost:5555/v1"
API_KEY="your-api-key-here"

echo "🚀 Testing JWT Authentication API"
echo "=================================="
echo

# Step 1: Create a test user
echo "📝 Step 1: Creating a test user..."
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/user" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: ${API_KEY}" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "Test@123456"
  }')

echo "Create User Response:"
echo "$CREATE_RESPONSE" | jq '.' 2>/dev/null || echo "$CREATE_RESPONSE"
echo

# Step 2: Authenticate the user
echo "🔐 Step 2: Authenticating user..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: ${API_KEY}" \
  -d '{
    "email": "test@example.com",
    "password": "Test@123456"
  }')

echo "Login Response:"
echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"

# Extract JWT token from response
JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token' 2>/dev/null)

if [ "$JWT_TOKEN" != "null" ] && [ "$JWT_TOKEN" != "" ]; then
    echo
    echo "✅ JWT Token obtained: ${JWT_TOKEN:0:50}..."
    echo

    # Step 3: Access protected route with JWT token
    echo "👤 Step 3: Accessing protected profile endpoint..."
    PROFILE_RESPONSE=$(curl -s -X GET "${BASE_URL}/user/profile" \
      -H "Authorization: Bearer ${JWT_TOKEN}" \
      -H "X-API-Key: ${API_KEY}")

    echo "Profile Response:"
    echo "$PROFILE_RESPONSE" | jq '.' 2>/dev/null || echo "$PROFILE_RESPONSE"
    echo

    # Step 4: Test with invalid token
    echo "❌ Step 4: Testing with invalid token..."
    INVALID_RESPONSE=$(curl -s -X GET "${BASE_URL}/user/profile" \
      -H "Authorization: Bearer invalid-token" \
      -H "X-API-Key: ${API_KEY}")

    echo "Invalid Token Response:"
    echo "$INVALID_RESPONSE" | jq '.' 2>/dev/null || echo "$INVALID_RESPONSE"
else
    echo "❌ Failed to obtain JWT token. Check the login response above."
fi

echo
echo "🎉 Test completed!"
echo
echo "To run this test:"
echo "1. Make sure your server is running: go run cmd/api/main.go"
echo "2. Make sure MongoDB is running"
echo "3. Set the correct API_KEY in the script"
echo "4. Run: chmod +x test_auth.sh && ./test_auth.sh"
