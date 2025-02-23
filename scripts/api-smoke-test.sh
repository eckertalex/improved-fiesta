#!/bin/bash

API_URL="http://localhost:45067/v1"
AUTH_URL="$API_URL/tokens/session"
EMAIL="admin@improved-fiesta.go"
PASSWORD="admin123"

echo_step() {
	echo -e "\n\033[1;34m$1\033[0m"
}

request() {
	local method="$1"
	local endpoint="$2"
	local data="$3"
	local expected_error="$4"
	local result

	if [ -n "$data" ]; then
		result=$(curl -s -H "Authorization: Bearer $token" -X "$method" -d "$data" "$API_URL/$endpoint")
	else
		result=$(curl -s -H "Authorization: Bearer $token" -X "$method" "$API_URL/$endpoint")
	fi

	if [ -n "$expected_error" ]; then
		if echo "$result" | jq -e ".error | contains(\"$expected_error\")" >/dev/null 2>/dev/null; then
			echo "✅ Passed"
		elif echo "$result" | jq -e ".error[] | contains(\"$expected_error\")" >/dev/null 2>/dev/null; then
			echo "✅ Passed"
		else
			echo "❌ Failed: Expected error '$expected_error' but got: $result" >&2
		fi
	else
		if [ -z "$result" ] || echo "$result" | jq -e ".data" >/dev/null 2>/dev/null; then
			echo "✅ Passed"
		else
			echo "❌ Failed: Unexpected response: $result" >&2
		fi
	fi
}

echo_step "Checking if API is reachable..."
if ! curl -s --head "$API_URL" >/dev/null; then
	echo "Error: API is not reachable at $API_URL" >&2
	exit 1
fi

echo_step "Obtaining authentication token..."
token=$(curl -s -X POST -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" "$AUTH_URL" | jq -r ".authentication_token.token")
if [ -z "$token" ] || [ "$token" == "null" ]; then
	echo "Error: Failed to obtain authentication token" >&2
	exit 1
fi

echo_step "Should fail with 'cannot remove the last admin user'"
request PATCH "users/1/role" "{\"role\": \"user\"}" "cannot remove the last admin user"

echo_step "Should fail with 'cannot remove the last admin user'"
request DELETE "users/1" "" "cannot remove the last admin user"

echo_step "Should pass changing username"
request PATCH "users/1" "{\"username\": \"nimda\"}"

echo_step "Should fail with 'a user with this username already exists'"
request PATCH "users/2" "{\"username\": \"nimda\"}" "a user with this username already exists"

echo_step "Should fail with 'a user with this email already exists'"
request PATCH "users/2" "{\"email\": \"admin@improved-fiesta.go\"}" "a user with this email address already exists"

echo_step "Should pass getting user details"
request GET "users/1"

echo_step "Should pass getting another user details"
request GET "users/2"

echo_step "Should pass getting all users"
request GET "users"

echo_step "Should pass getting all users with pagination"
request GET "users?page=1&page_size=2&sort=-id"

echo_step "Should pass searching for partial username"
request GET "users?username=nim"

echo_step "Should pass getting debug metrics"
metrics=$(curl -s -H "Authorization: Bearer $token" "http://localhost:45067/debug/vars" | jq -r ".version")
echo "$metrics"

echo_step "Should pass logging out"
request DELETE "tokens/session"

echo_step "Should fail with 'invalid or missing authentication token'"
request GET "users/1" "" "invalid or missing authentication token"
