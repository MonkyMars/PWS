# Custom Hooks

This directory contains custom React hooks for data fetching, state management, and business logic.

## Available Hooks

### Authentication Hooks (`use-auth.ts`)

- `useLogin()` - Handles user authentication
- `useRegister()` - Handles user registration
- `useLogout()` - Handles user logout
- `useCurrentUser()` - Fetches current user data
- `useIsAuthenticated()` - Returns authentication status

### Subject Hooks (`use-subjects.ts`)

- `useSubjects(filters?)` - Fetches user's subjects with optional filters
- `useSubject(subjectId)` - Fetches specific subject details
- `useAnnouncements(filters?)` - Fetches announcements with pagination
- `useSubjectFiles(filters?)` - Fetches subject files with pagination
- `useUploadFile()` - Handles file upload with progress tracking
- `useDeleteFile()` - Handles file deletion

## Features

### Automatic Caching

All hooks use TanStack Query for intelligent caching and background updates.

### Error Handling

Comprehensive error handling with user-friendly error messages in Dutch.

### Optimistic Updates

Mutations automatically update the cache for better UX.

### Type Safety

Full TypeScript support with proper return types and generics.

## Usage Examples

```typescript
// Get current user
const { data: user, isLoading, error } = useCurrentUser();

// Get subjects with loading state
const { data: subjects, isLoading } = useSubjects();

// Upload file with progress
const uploadMutation = useUploadFile();
uploadMutation.mutate({
  file,
  subjectId,
  onProgress: (progress) => console.log(`${progress}%`),
});
```

## Best Practices

1. Always handle loading and error states in components
2. Use the `enabled` option to conditionally fetch data
3. Implement proper cleanup in components using these hooks
4. Use optimistic updates for better perceived performance
