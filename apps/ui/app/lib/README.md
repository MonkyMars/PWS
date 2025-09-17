# Library Utilities

This directory contains utility libraries, configurations, and helper functions.

## Files

### `api-client.ts`

HTTP client for API communication with the backend server.

**Features:**

- Automatic token management
- Request/response interceptors
- Error handling
- File upload support with progress tracking
- TypeScript-safe API responses

**Usage:**

```typescript
import { apiClient } from "~/lib/api-client";

// GET request
const response = await apiClient.get<User[]>("/users");

// POST request
const response = await apiClient.post<User>("/users", userData);

// File upload
const response = await apiClient.uploadFile(
  "/files/upload",
  file,
  { subjectId },
  (progress) => console.log(`${progress}%`)
);
```

### `query-client.ts`

TanStack Query client configuration for data fetching and caching.

**Features:**

- Optimized caching strategy
- Automatic retry logic
- Global error handling
- Performance-tuned defaults

**Configuration:**

- Stale time: 5 minutes
- Cache time: 10 minutes
- Retry attempts: 2 for queries, 1 for mutations
- No refetch on window focus

## Design Principles

### Type Safety

All utilities are fully typed with TypeScript, providing compile-time safety and excellent developer experience.

### Performance

Configurations are optimized for the educational environment use case, balancing data freshness with performance.

### Error Handling

Comprehensive error handling with fallbacks and user-friendly error messages in Dutch.

### Extensibility

Modular design allows for easy extension and customization of functionality.

## Environment Variables

The API client reads configuration from environment variables:

- `API_URL` - Backend API base URL (default: http://localhost:8080/api)
