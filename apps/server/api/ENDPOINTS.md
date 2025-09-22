# Endpoints

This is the documentation for all the endpoints and how they work.

## Available Endpoints
- GET /health - Returns server health plus some metrics like go routines and memory usage.
- GET /health/database - Returns database connection status and the latency
- GET /* - Fallback route, returns 404 