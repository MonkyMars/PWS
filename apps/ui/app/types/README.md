# TypeScript Types

This directory contains TypeScript type definitions for the PWS ELO application.

## Type Categories

### Authentication Types (`auth.ts`)

- `UserRole` - User role enumeration (student, teacher, admin)
- `User` - Complete user interface
- `LoginCredentials` - Login form data
- `RegisterData` - Registration form data
- `AuthResponse` - API authentication response
- `AuthState` - Client-side authentication state

### Subject Types (`subject.ts`)

- `Subject` - Academic subject interface
- `Announcement` - Subject announcement interface
- `SubjectFile` - File attachment interface
- `SubjectWithDetails` - Extended subject with related data
- `SubjectEnrollment` - User-subject relationship

### API Types (`api.ts`)

- `ApiResponse<T>` - Standard API response wrapper
- `ApiError` - Error response structure
- `PaginationParams` - Query parameters for pagination
- `PaginatedResponse<T>` - Paginated data response
- `UploadProgress` - File upload progress tracking
- Various filter interfaces for data querying

## Design Principles

### Type Safety

All types are designed to prevent runtime errors and provide excellent IDE support with:

- Strict null checks
- Exhaustive union type checking
- Generic type parameters where appropriate
- No use of `any` type

### Consistency

Types follow consistent naming conventions:

- Interfaces use PascalCase
- Properties use camelCase
- Enums use camelCase values
- Generic types use descriptive single letters (T, U, K, V)

### Extensibility

Types are designed for future growth:

- Interfaces can be extended without breaking changes
- Union types allow for easy addition of new variants
- Generic types provide flexibility for different data shapes

## Usage Guidelines

### Importing Types

```typescript
// Import specific types
import type { User, Subject } from "~/types";

// Import all types (avoid in large files)
import type * as Types from "~/types";
```

### Extending Types

```typescript
// Extend existing interfaces
interface ExtendedUser extends User {
  customField: string;
}

// Create type unions
type UserOrGuest = User | { type: "guest" };
```

### Generic Usage

```typescript
// Use with API responses
const response: ApiResponse<User[]> = await apiClient.get("/users");

// Use with pagination
const paginatedUsers: PaginatedResponse<User> =
  await apiClient.get("/users/paginated");
```

## Validation

Types are complemented by runtime validation using Zod schemas in the authentication components. This ensures type safety both at compile time and runtime.
