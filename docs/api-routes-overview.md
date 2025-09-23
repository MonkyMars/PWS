# API Routes Overview

This document provides a comprehensive overview of all API routes in the PWS application, organized by functionality and showing the complete routing structure.

## Route Structure

The API uses a modular routing structure with grouped endpoints for better organization and maintainability.

### Base URL
```
http://localhost:8082  (development)
https://yourdomain.com (production)
```

## Authentication Routes (`/auth`)

### Basic Authentication
- **POST** `/auth/login`
  - **Description**: Authenticate user with email/password
  - **Auth Required**: No
  - **Request Body**: `{ "email": "string", "password": "string" }`
  - **Response**: User data with JWT tokens in secure cookies
  - **Status Codes**: 200 (success), 400 (bad request), 401 (unauthorized)

- **POST** `/auth/register`
  - **Description**: Register new user account
  - **Auth Required**: No
  - **Request Body**: `{ "username": "string", "email": "string", "password": "string" }`
  - **Response**: User data with JWT tokens
  - **Status Codes**: 200 (success), 400 (validation error), 409 (user exists)

- **POST** `/auth/refresh`
  - **Description**: Refresh access token using refresh token from cookies
  - **Auth Required**: No (uses refresh token cookie)
  - **Response**: New access and refresh tokens in cookies
  - **Status Codes**: 200 (success), 401 (invalid/expired token)

- **POST** `/auth/logout`
  - **Description**: Logout user, blacklist tokens, clear cookies
  - **Auth Required**: No (handles invalid tokens gracefully)
  - **Response**: Success message
  - **Status Codes**: 200 (success)

- **GET** `/auth/me`
  - **Description**: Get current authenticated user information
  - **Auth Required**: Yes (Access Token)
  - **Response**: Current user data
  - **Status Codes**: 200 (success), 401 (unauthorized), 404 (user not found)

### Google OAuth Routes (`/auth/google`)

- **GET** `/auth/google/url`
  - **Description**: Generate Google OAuth authorization URL
  - **Auth Required**: Yes (Access Token)
  - **Response**: `{ "data": "https://accounts.google.com/oauth/authorize?..." }`
  - **Status Codes**: 200 (success), 401 (unauthorized), 500 (state generation failed)

- **GET** `/auth/google/callback`
  - **Description**: Handle Google OAuth callback (automatic redirect)
  - **Auth Required**: No (validates state token)
  - **Query Parameters**: `state`, `code` (provided by Google)
  - **Response**: Redirect to frontend success page
  - **Status Codes**: 302 (redirect), 400 (invalid callback)

- **GET** `/auth/google/access-token`
  - **Description**: Get fresh Google Drive API access token
  - **Auth Required**: Yes (Access Token)
  - **Response**: 
    ```json
    {
      "data": {
        "access_token": "ya29.xxx",
        "expiry": "2024-01-01T12:00:00Z",
        "token_type": "Bearer"
      }
    }
    ```
  - **Status Codes**: 200 (success), 401 (unauthorized/no linked account), 500 (refresh failed)

- **GET** `/auth/google/status`
  - **Description**: Check if user has linked their Google account
  - **Auth Required**: Yes (Access Token)
  - **Response**: `{ "data": { "linked": true, "user_id": "uuid" } }`
  - **Status Codes**: 200 (success), 401 (unauthorized), 500 (database error)

- **DELETE** `/auth/google/unlink`
  - **Description**: Unlink user's Google account
  - **Auth Required**: Yes (Access Token)
  - **Response**: `{ "data": { "message": "Google account unlinked successfully" } }`
  - **Status Codes**: 200 (success), 401 (unauthorized), 500 (database error)

## File Routes (`/files`)

*Note: File routes implementation details would be documented here when available*

## Health Check Routes

- **GET** `/health`
  - **Description**: Get server health status with metrics
  - **Auth Required**: No
  - **Response**: Server health data (memory, goroutines, etc.)
  - **Status Codes**: 200 (healthy), 503 (unhealthy)

- **GET** `/health/database`
  - **Description**: Check database connection status and latency
  - **Auth Required**: No
  - **Response**: Database connection status and response time
  - **Status Codes**: 200 (connected), 503 (connection failed)

## Fallback Route

- **GET** `/*` (catch-all)
  - **Description**: Fallback route for undefined endpoints
  - **Response**: 404 Not Found
  - **Status Codes**: 404

## Authentication Flow

### Standard Authentication Flow
1. **Register/Login** → Get JWT tokens in secure cookies
2. **Use Access Token** → Make authenticated requests
3. **Token Refresh** → Automatically refresh when access token expires
4. **Logout** → Clear tokens and blacklist them

### Google OAuth Flow
1. **Must be logged in** → User needs valid access token first
2. **Get OAuth URL** → `GET /auth/google/url`
3. **User Authorization** → Redirect to Google consent screen
4. **OAuth Callback** → Google redirects to `/auth/google/callback`
5. **Link Complete** → User redirected to frontend success page
6. **Use Google APIs** → Get access tokens via `/auth/google/access-token`

## Request/Response Format

### Standard Response Format
```json
{
  "success": true,
  "data": "response_data_here",
  "message": "optional_message"
}
```

### Error Response Format
```json
{
  "success": false,
  "error": "error_message",
  "details": "optional_detailed_error_info"
}
```

### Validation Error Format
```json
{
  "success": false,
  "error": "Validation failed",
  "validation_errors": [
    {
      "field": "email",
      "message": "Email is required",
      "value": ""
    }
  ]
}
```

## Security Features

### Authentication
- **JWT Tokens**: Access tokens (15 min) and refresh tokens (7 days)
- **Secure Cookies**: HTTPOnly, Secure, SameSite cookies
- **Token Rotation**: Refresh tokens are rotated on each use
- **Token Blacklisting**: Revoked tokens stored in Redis

### Google OAuth Security
- **CSRF Protection**: Cryptographically secure state tokens
- **Single-Use States**: State tokens expire after 10 minutes and are deleted after use
- **Minimal Scopes**: Read-only Google Drive access by default
- **Server-Side Storage**: Refresh tokens stored securely server-side

### Rate Limiting & CORS
- **CORS**: Configured for cross-origin requests
- **Request Logging**: All requests logged with correlation IDs
- **Error Handling**: Graceful error handling with proper status codes

## Environment Configuration

Routes behavior can be configured via environment variables:

```bash
# Google OAuth
GOOGLE_OAUTH_CLIENT_ID=your_client_id
GOOGLE_OAUTH_CLIENT_SECRET=your_client_secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8082/auth/google/callback
FRONTEND_URL=http://localhost:3000

# Server
PORT=8082
LOG_LEVEL=info

# Cache (Redis) - Required for OAuth state management
CACHE_ADDRESS=localhost:6379
```

## Development vs Production

### Development
- Detailed error messages
- Debug logging available
- HTTP localhost URLs supported
- Test users allowed for unverified OAuth apps

### Production
- Generic error messages for security
- Structured logging
- HTTPS required for OAuth redirects
- OAuth consent screen must be verified by Google

## Monitoring & Debugging

### Health Endpoints
Use `/health` and `/health/database` for monitoring server and database status.

### Logging
- Request/response logging with correlation IDs
- Error logging with context
- OAuth flow logging for debugging
- Set `LOG_LEVEL=debug` for detailed logs

### Common Issues
- **OAuth "redirect_uri_mismatch"**: Check `GOOGLE_OAUTH_REDIRECT_URL` matches Google Console
- **"Invalid state"**: Check Redis connection and system time sync
- **"No refresh token"**: User may need to re-consent with `prompt=consent`

This routing structure provides a clean, secure, and scalable API architecture for the PWS application.