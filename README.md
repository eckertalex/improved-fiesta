# Improved Fiesta 🎉

A Go-based API service with SQLite storage and user authentication.

## Getting Started

### Prerequisites

- Go 1.23.2
- SQLite3
- Make

### Quick Start

1. Clone the repository
2. Run `make tidy` to install dependencies
3. Copy `.env.example` to `.env`:

```sh
cp .env.example .env
```

4. Run database migrations to set up the database schema:

```sh
make db/migrations/up
```

5. Seed the database:

```sh
make db/seed
```

6. Start the development server with hot reload:

```sh
make watch/api
```

The API will be available at the port specified in your `.env` file (default: 45067).

### Available Make Commands

Run `make help` to see all available commands. Here are the most commonly used ones:

#### Development

- `make run/api` - Build and run the API
- `make watch/api` - Run the API with hot reload using Air
- `make tidy` - Format code and tidy up Go modules

#### Testing

- `make test` - Run all tests
- `make itest` - Run integration tests
- `make audit` - Run quality control checks (tests, format verification, vulnerability scanning)

#### Database Operations

- `make db/migrations/new name=migration_name` - Create a new migration
- `make db/migrations/up` - Apply all pending migrations
- `make db/migrations/down` - Revert all migrations
- `make db/seed` - Seed the database with initial data
- `make db/connect` - Connect to the SQLite database

## API Documentation

This document provides comprehensive documentation for all API endpoints, including required permissions, request/response formats, and examples.

### Table of Contents

| Method | Endpoint                                                  | Permission     | Description               |
| ------ | --------------------------------------------------------- | -------------- | ------------------------- |
| GET    | [/v1/healthcheck](#get-v1healthcheck)                     | Public         | Application health status |
| GET    | [/v1/users](#get-v1users)                                 | Admin          | List users                |
| POST   | [/v1/users](#post-v1users)                                | Public         | Create new user           |
| GET    | [/v1/users/:id](#get-v1usersid)                           | Admin or Owner | Get user details          |
| PATCH  | [/v1/users/:id](#patch-v1usersid)                         | Admin or Owner | Update user details       |
| DELETE | [/v1/users/:id](#delete-v1usersid)                        | Admin or Owner | Delete user account       |
| PATCH  | [/v1/users/:id/role](#patch-v1usersid-role)               | Admin          | Update user role          |
| POST   | [/v1/users/activate](#post-v1usersactivate)               | Public         | Activate user account     |
| POST   | [/v1/users/reset-password](#post-v1usersreset-password)   | Public         | Reset user password       |
| POST   | [/v1/tokens/session](#post-v1tokenssession)               | Public         | Generate auth token       |
| DELETE | [/v1/tokens/session](#delete-v1tokenssession)             | Public         | Delete auth token         |
| POST   | [/v1/tokens/activation](#post-v1tokensactivation)         | Public         | Request activation token  |
| POST   | [/v1/tokens/password-reset](#post-v1tokenspassword-reset) | Public         | Request password reset    |
| GET    | [/debug/vars](#get-debugvars)                             | Admin          | Debug metrics             |

### Authentication

For endpoints requiring authentication, include the authentication token in the Authorization header:

```
Authorization: Bearer <token>
```

Tokens can be obtained using the [/v1/auth/authentication](#post-v1tokensauthentication) endpoint.

### Endpoint Details

#### GET /v1/healthcheck

**Permission:** Public  
**Description:** Check application health status and version information

**Example Request:**

```bash
curl http://localhost:45067/v1/healthcheck
```

**Response Body:**

- `status` (string): Current application status
- `system_info` (object):
  - `environment` (string): Current environment
  - `version` (string): Application version

**Example Response:**

```json
{
  "status": "available",
  "system_info": {
    "environment": "development",
    "version": "2024-01-15T10:30:00Z-c3418de38b57f0a8bf49d9f7759469eac37cb410"
  }
}
```

#### POST /v1/users

**Permission:** Public  
**Description:** Create a new user account

**Request Body:**

- `username` (string): Username
- `email` (string): User's email address
- `password` (string): User's password

**Example Request:**

```bash
curl -X POST \
     -d '{"username": "bob", "email": "bob.smith@example.com", "password": "SecurePass123!"}' \
     http://localhost:45067/v1/users
```

**Response Body:**

- `data` (object): User details

**Example Response:**

```json
{
  "data": {
    "id": 124,
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z",
    "username": "bob",
    "email": "bob.smith@example.com",
    "activated": false,
    "role": "user"
  }
}
```

#### GET /v1/users/:id

**Permission:** Admin or Owner  
**Description:** Retrieve user details by ID

**Path Parameters:**

- `id` (number): User ID

**Request Headers:**

- `Authorization`: Bearer token

**Example Request:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:45067/v1/users/124
```

**Response Body:**

- `data` (object): User details object

**Example Response:**

```json
{
  "data": {
    "id": 124,
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z",
    "username": "bob",
    "email": "bob.smith@example.com",
    "activated": false,
    "role": "user"
  }
}
```

#### PATCH /v1/users/:id

**Permission:** Admin or Owner  
**Description:** Update user details

**Path Parameters:**

- `id` (number): User ID

**Request Headers:**

- `Authorization`: Bearer token

**Request Body:** (all fields optional)

- `username` (string): New username
- `email` (string): New email
- `password` (string): New password

**Example Request:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
     -X PATCH \
     -d '{"username": "robert"}' \
     http://localhost:45067/v1/users/124
```

**Response Body:**

- `data` (object): Updated user details

**Example Response:**

```json
{
  "data": {
    "id": 124,
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T15:30:00Z",
    "username": "robert",
    "email": "bob.smith@example.com",
    "activated": false,
    "role": "user"
  }
}
```

#### DELETE /v1/users/:id

**Permission:** Admin or Owner  
**Description:** Delete user account

**Path Parameters:**

- `id` (number): User ID

**Request Headers:**

- `Authorization`: Bearer token

**Example Request:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
     -X DELETE \
     http://localhost:45067/v1/users/124
```

**Response:** 204 No Content

#### PATCH /v1/users/:id/role

**Permission:** Admin  
**Description:** Update user's role

**Path Parameters:**

- `id` (number): User ID

**Request Headers:**

- `Authorization`: Bearer token

**Request Body:**

- `role` (string): New role (must be 'admin' or 'user')

**Example Request:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
     -X PATCH \
     -d '{"role": "admin"}' \
     http://localhost:45067/v1/users/124/role
```

**Response Body:**

- `data` (object): Updated user details

**Example Response:**

```json
{
  "data": {
    "id": 124,
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T16:00:00Z",
    "username": "robert",
    "email": "bob.smith@example.com",
    "activated": false,
    "role": "admin"
  }
}
```

#### POST /v1/users/activate

**Permission:** Public  
**Description:** Activate user account

**Request Body:**

- `token` (string): Activation token

**Example Request:**

```bash
curl -X POST \
     -d '{"token": "URQXUFRJDA4KCIZ3CLZTHS3K3U"}' \
     http://localhost:45067/v1/users/activate
```

**Response Body:**

- `data` (object): Updated user details

**Example Response:**

```json
{
  "data": {
    "id": 123,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "username": "alice",
    "email": "alice.johnson@example.com",
    "activated": true,
    "role": "user"
  }
}
```

#### POST /v1/users/reset-password

**Permission:** Public  
**Description:** Reset user password

**Request Body:**

- `password` (string): New password
- `token` (string): Password reset token

**Example Request:**

```bash
curl -X POST \
     -d '{"password": "SecurePass123!", "token": "K5YB7CA3KTZL4B6YXO5JLNLWWU"}' \
     http://localhost:45067/v1/users/reset-password
```

**Response Body:**

- `data` (object): Updated user details

**Example Response:**

```json
{
  "data": {
    "id": 123,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z",
    "username": "alice",
    "email": "alice.johnson@example.com",
    "activated": true,
    "role": "user"
  }
}
```

#### POST /v1/tokens/session

**Permission:** Public  
**Description:** Generate authentication token

**Request Body:**

- `email` (string): User's email
- `password` (string): User's password

**Example Request:**

```bash
curl -X POST \
     -d '{"email": "admin@example.com", "password": "AdminPass123!"}' \
     http://localhost:45067/v1/tokens/session
```

**Response Body:**

- `authentication_token` (object):
  - `token` (string): Authentication token
  - `expiry` (string): Token expiration timestamp

**Example Response:**

```json
{
  "authentication_token": {
    "token": "QNSMBHTILTB4RDRTMRYACTGCTE",
    "expiry": "2024-01-15T13:30:00Z"
  }
}
```

#### POST /v1/tokens/activation

**Permission:** Public  
**Description:** Request new activation token

**Request Body:**

- `email` (string): User's email

**Example Request:**

```bash
curl -X POST \
     -d '{"email": "alice.johnson@example.com"}' \
     http://localhost:45067/v1/tokens/activation
```

**Response Body:**

- `message` (string): Confirmation message

**Example Response:**

```json
{
  "message": "an email will be sent to you containing activation instructions"
}
```

#### POST /v1/tokens/password-reset

**Permission:** Public  
**Description:** Request password reset token

**Request Body:**

- `email` (string): User's email

**Example Request:**

```bash
curl -X POST \
     -d '{"email": "alice.johnson@example.com"}' \
     http://localhost:45067/v1/tokens/password-reset
```

**Response Body:**

- `message` (string): Confirmation message

**Example Response:**

```json
{
  "message": "an email will be sent to you containing password reset instructions"
}
```

#### GET /debug/vars

**Permission:** Admin  
**Description:** Displays application metrics for debugging purposes

**Request Headers:**

- `Authorization`: Bearer token

**Example Request:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:45067/debug/vars
```

**Response Body:**

- Various debugging metrics and application variables

## Todo

- [ ] Add GET /v1/users endpoint
- [ ] Write End-to-End tests
