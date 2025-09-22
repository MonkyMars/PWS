# Endpoints

This is the documentation for all the endpoints and how they work.

## Available Endpoints

### General Endpoints
- GET /health - Returns server health plus some metrics like go routines and memory usage.
- GET /health/database - Returns database connection status and the latency
- GET /* - Fallback route, returns 404

### Auth Endpoints
- POST /auth/login - Login user and return JWT tokens in cookies
- POST /auth/register - Register new user and return JWT tokens in cookies
- POST /auth/refresh - Refresh access token using refresh token in cookies
- POST /auth/logout - Logout user, blacklist tokens and clear cookies
- GET /auth/me - Get current authenticated user info (requires valid access token)
