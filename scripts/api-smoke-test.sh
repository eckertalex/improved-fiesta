#!/bin/bash

API_URL="http://localhost:45067/v1"
AUTH_URL="$API_URL/tokens/session"
EMAIL="admin@improved-fiesta.go"
PASSWORD="admin123"

if ! curl -s --head "$API_URL" > /dev/null; then
    echo "Error: API is not reachable at $API_URL" >&2
    exit 1
fi

token=$(curl -s -X POST -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" "$AUTH_URL" | jq -r ".authentication_token.token")
if [ -z "$token" ]; then
	echo "Error: Failed to obtain authentication token"
	exit 1
fi

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"role\": \"user\"}" "$API_URL/users/1/role"

curl -s -H "Authorization: Bearer $token" -X DELETE "$API_URL/users/1"

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"username\": \"nimda\"}" "$API_URL/users/1"

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"username\": \"nimda\"}" "$API_URL/users/2"

curl -s -H "Authorization: Bearer $token" -X PATCH -d "{\"email\": \"admin@improved-fiesta.go\"}" "$API_URL/users/2"

curl -s -H "Authorization: Bearer $token" "$API_URL/users/1"

curl -s -H "Authorization: Bearer $token" "$API_URL/users/2"

curl -s -H "Authorization: Bearer $token" "$API_URL/users" | jq | pbcopy

curl -s -H "Authorization: Bearer $token" "$API_URL/users?page=1&page_size=2&sort=-id"

curl -s -H "Authorization: Bearer $token" "$API_URL/users?username=nim"

curl -s -H "Authorization: Bearer $token" "http://localhost:45067/debug/vars" | jq -r ".version"

curl -s -H "Authorization: Bearer $token" -X DELETE "$API_URL/tokens/session"

curl -s -H "Authorization: Bearer $token" "$API_URL/users/1"
