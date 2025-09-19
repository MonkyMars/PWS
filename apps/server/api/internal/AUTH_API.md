# Authentication API Documentation

This document describes the authentication endpoints available in the PWS application.

## Base URL
```
http://localhost:8082/auth
```

## Endpoints

### 1. User Registration
**POST** `/auth/register`

Registers a new user account and returns authentication tokens.

#### Request Body
```json
{
  "username": "string",
  "email": "string", 
  "password": "string"
}
```

#### Validation Rules
- `username`: Required, non-empty string
- `email`: Required, non-empty string, must be unique
- `password`: Required, minimum 6 characters

#### Success Response (201)
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "role": "student"
    },
    "access_token": "jwt_string",
    "refresh_token": "jwt_string"
  },
  "message": "Success"
}
```

#### Error Responses
- **400 Bad Request**: Invalid request body or validation errors
- **409 Conflict**: User with email or username already exists
- **500 Internal Server Error**: Server error during registration

---

### 2. User Login
**POST** `/auth/login`

Authenticates a user and returns authentication tokens.

#### Request Body
```json
{
  "email": "string",
  "password": "string"
}
```

#### Validation Rules
- `email`: Required, non-empty string
- `password`: Required, non-empty string

#### Success Response (200)
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "username": "string", 
      "email": "string",
      "role": "student"
    },
    "access_token": "jwt_string",
    "refresh_token": "jwt_string"
  },
  "message": "Success"
}
```

#### Error Responses
- **400 Bad Request**: Invalid request body or validation errors
- **401 Unauthorized**: Invalid email or password
- **500 Internal Server Error**: Server error during login

---

### 3. Token Refresh
**POST** `/auth/refresh`

Refreshes an expired access token using a valid refresh token.

#### Request Body
```json
{
  "refresh_token": "jwt_string"
}
```

#### Validation Rules
- `refresh_token`: Required, valid JWT refresh token

#### Success Response (200)
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string", 
      "role": "student"
    },
    "access_token": "jwt_string",
    "refresh_token": "jwt_string"
  },
  "message": "Success"
}
```

#### Error Responses
- **400 Bad Request**: Invalid request body
- **401 Unauthorized**: Invalid or expired refresh token
- **500 Internal Server Error**: Server error during token refresh

---

### 4. Get Current User
**GET** `/auth/me`

Returns the current authenticated user's information.

#### Headers
```
Authorization: Bearer <access_token>
```

#### Success Response (200)
```json
{
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "role": "student"
  },
  "message": "Success"
}
```

#### Error Responses
- **401 Unauthorized**: Missing, invalid, or expired access token
- **500 Internal Server Error**: Server error

---

### 5. User Logout
**POST** `/auth/logout`

Logs out the current user. In a complete implementation, this would invalidate the token.

#### Headers
```
Authorization: Bearer <access_token>
```

#### Success Response (200)
```json
{
  "data": {
    "message": "Logged out successfully"
  },
  "message": "Success"
}
```

#### Error Responses
- **401 Unauthorized**: Missing, invalid, or expired access token

---

## Authentication Flow

### Initial Authentication
1. User registers with `/auth/register` or logs in with `/auth/login`
2. Server returns both access and refresh tokens
3. Client stores both tokens securely

### Making Authenticated Requests
1. Include access token in Authorization header: `Bearer <access_token>`
2. If request returns 401, use refresh token to get new access token
3. Retry original request with new access token

### Token Refresh Flow
1. When access token expires (15 minutes default), use `/auth/refresh`
2. Send refresh token to get new access and refresh tokens
3. Update stored tokens and continue making requests

### Logout Flow
1. Call `/auth/logout` with valid access token
2. Clear stored tokens from client
3. Redirect to login page

## Token Configuration

### Access Token
- **Expiry**: 15 minutes (configurable via `ACCESS_TOKEN_EXPIRY`)
- **Purpose**: Authorizing API requests
- **Storage**: Memory/session storage (recommended)

### Refresh Token  
- **Expiry**: 7 days (configurable via `REFRESH_TOKEN_EXPIRY`)
- **Purpose**: Obtaining new access tokens
- **Storage**: Secure HTTP-only cookie (recommended) or localStorage

## Security Considerations

### Current Implementation
- JWT tokens signed with HMAC-SHA256
- Passwords hashed with bcrypt
- Basic input validation
- CORS configured for localhost:5173

## Error Response Format

All error responses follow this structure:

```json
{
  "error": "error_message",
  "code": "ERROR_CODE",
  "details": "additional_details"
}
```

### Validation Error Format
```json
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR", 
  "validation_errors": [
    {
      "field": "field_name",
      "message": "error_message",
      "value": "submitted_value"
    }
  ]
}
```

## Example Usage

### TypeScript
```typescript
// Login
const loginResponse = await fetch('/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

// Get current user
const userResponse = await fetch('/auth/me', {
  headers: { 
    'Authorization': `Bearer ${accessToken}`
  }
});

// Refresh token
const refreshResponse = await fetch('/auth/refresh', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    refresh_token: refreshToken
  })
});
```

### cURL Examples
```bash
# Register
curl -X POST http://localhost:8082/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login  
curl -X POST http://localhost:8082/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get current user
curl -X GET http://localhost:8082/auth/me \
  -H "Authorization: Bearer <access_token>"

# Refresh token
curl -X POST http://localhost:8082/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'

# Logout
curl -X POST http://localhost:8082/auth/logout \
  -H "Authorization: Bearer <access_token>"
```
