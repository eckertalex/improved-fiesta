#!/bin/bash

API_URL="http://localhost:45067/v1"
AUTH_URL="$API_URL/tokens/session"
EMAIL="admin@improved-fiesta.go"
PASSWORD="admin123"

token=$(curl -s -X POST -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" "$AUTH_URL" | jq -r ".authentication_token.token")
if [ -z "$token" ]; then
	echo "Failed to obtain authentication token"
	exit 1
fi

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"role\": \"user\"}" "$API_URL/users/1/role"

sleep 1

curl -s -H "Authorization: Bearer $token" -X DELETE "$API_URL/users/1"

sleep 1

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"name\": \"ADMIN\"}" "$API_URL/users/1"

sleep 1

curl -s -H "Authorization: Bearer $token" "$API_URL/users/1"

sleep 1

curl -s -H "Authorization: Bearer $token" "$API_URL/users/2"
