# Routes

This directory contains React Router v7 route components that define the application's pages.

## Available Routes

### Public Routes

- `home.tsx` (`/`) - Landing page explaining ELO concept and features
- `login.tsx` (`/login`) - User authentication page
- `register.tsx` (`/register`) - User registration page

### Protected Routes

- `dashboard.tsx` (`/dashboard`) - Main dashboard with subject overview
- `subjects.$subjectId.tsx` (`/subjects/:subjectId`) - Detailed subject view

## Route Structure

### File Naming Convention

React Router v7 uses file-based routing with specific naming conventions:

- `home.tsx` - Maps to `/`
- `login.tsx` - Maps to `/login`
- `subjects.$subjectId.tsx` - Maps to `/subjects/:subjectId` (dynamic route)

### Route Components

Each route component includes:

- `meta()` function for page metadata (title, description)
- Default export with the page component
- Authentication checks for protected routes
- Loading states and error handling

## Authentication Flow

### Protected Routes

Protected routes (dashboard, subject details) check for authentication:

1. Show loading spinner while checking auth status
2. Redirect to `/login` if not authenticated
3. Render the page content if authenticated

### Public Routes

Public routes (home, login, register) are accessible to all users but may show different content based on authentication state.

## Navigation

Routes are connected through:

- `<Link>` components for client-side navigation
- `<Navigate>` components for programmatic redirects
- Navigation component with conditional menu items

## Error Handling

Each route includes proper error handling:

- Loading states during data fetching
- Error boundaries for unexpected errors
- Fallback UI for missing data
- Redirect to appropriate pages when needed

## Performance

Routes are optimized for performance:

- Code splitting at the route level
- Lazy loading of heavy components
- Prefetching of likely next routes
- Optimized bundle sizes
