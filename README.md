# PWS (Profiel werkstuk)

A modern, full-stack web application built with Go backend and React frontend. The application features a robust API server with comprehensive response handling, configuration management, and a responsive user interface.

## Architecture Overview

PWS follows a monorepo structure with clearly separated frontend and backend applications:

```
PWS/
├── apps/
│   ├── server/          # Go backend API server
│   │   ├── main.go      # Application entry point
│   │   ├── api/         # API layer (routes, middleware, responses)
│   │   ├── config/      # Configuration management
│   │   ├── database/    # Database connection and models
│   │   ├── lib/         # Helper functions
│   │   ├── services/    # Business logic services
│   │   └── types/       # Shared type definitions
│   └── ui/              # React frontend application
│       ├── app/         # React Router v7 application
│       ├── public/      # Static assets
│       └── build/       # Production build output
├── LICENSE              # Project license
└── README.md           # This documentation
```

## Technology Stack

### Backend (Go Server)

- **Go 1.25.0** - High-performance programming language
- **Fiber v3.0.0-rc.1** - Express-inspired web framework
- **slog** - Structured logging (Go standard library)
- **goccy/go-json** - High-performance JSON encoding
- **godotenv** - Environment variable management

### Frontend (React UI)

- **React 19** - Modern UI library with latest features
- **React Router v7** - File-based routing with SSR support
- **TypeScript** - Type safety and developer experience
- **Tailwind CSS v4** - Utility-first CSS framework
- **Vite** - Fast build tool and development server
- **Bun** - Package manager and runtime

## Quick Start

### Prerequisites

- **Go 1.25.0** - [Download Go](https://golang.org/dl/)
- **Bun** (recommended) or **Node.js 18+** - [Install Bun](https://bun.sh/)
- **Git** - For version control

### Development Setup

1. **Clone the repository**:

   ```bash
   git clone git@github.com:MonkyMars/PWS.git
   cd PWS
   ```

2. **Start the backend server**:

   ```bash
   cd apps/server

   # Install dependencies
   go mod tidy

   # Copy environment configuration
   cp .env.example .env

   # Start development server
   go run main.go
   ```

   The API server will be available at `http://localhost:8080`

3. **Start the frontend application** (in a new terminal):

   ```bash
   cd apps/ui

   # Install bun
   npm install -g bun

   # Install dependencies
   bun install

   # Start development server
   bun run dev
   ```

   The UI will be available at `http://localhost:5173`

### Production Build

**Backend**:

```bash
cd apps/server
go build -o server main.go
./server
```

**Frontend**:

```bash
cd apps/ui
bun run build
bun run start
```

## Application Components

### Backend Architecture

#### Configuration Management (`/config/`)

- **Environment Variables**: Centralized loading and validation
- **Logging Setup**: Structured logging configuration with slog
- **Fiber Configuration**: Web server setup with optimized settings
- **Singleton Pattern**: Thread-safe access to configuration

#### Response Handling (`api/response/`)

- **Standardized Responses**: Consistent API response format
- **Builder Pattern**: Fluent API for response construction
- **Error Handling**: Structured error responses with codes and details
- **Pagination Support**: Built-in pagination for list endpoints

#### Middleware (`api/middleware/`)

- **CORS**: Cross-origin resource sharing configuration
- **Extensible Design**: Easy addition of new middleware components
- **Request/Response Processing**: Standardized middleware patterns

#### Routes (`api/routes/`)

- **RESTful Design**: Standard REST API patterns
- **Route Organization**: Feature-based route grouping
- **Handler Functions**: Clean separation of route logic

#### Internal (`api/internal/`)

- **Logic**: Contains all the business logic in an organized way
- **Core**: The core of the server
- **Health**: Clean health routes logic

### Frontend Architecture

#### Component Structure (`app/`)

- **Root Layout**: Base HTML structure and global providers
- **Route Components**: File-based routing with React Router v7
- **Welcome Component**: Landing page with responsive design
- **TypeScript Integration**: Full type safety throughout

#### Styling (`app/app.css`)

- **Tailwind CSS**: Utility-first styling approach
- **Responsive Design**: Mobile-first breakpoints
- **Component Styling**: Scoped styles with CSS modules support

## API Reference

### Response Format

All API endpoints return paginated responses in the following format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data
  },
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "pages": 10
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Responses

```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": {
      // Additional error context
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Health Check Endpoints

These endpoints are only available in the development environment

- `GET /health` - Basic server health check
- `GET /health/database` - Ping the database and check for errors

## Development Guidelines

### Backend Development

1. **Follow Go Conventions**: Use standard Go naming and structure conventions
2. **Document Public APIs**: Include Go doc comments for all exported functions
3. **Handle Errors Properly**: Use structured error responses with context
4. **Write Tests**: Comprehensive test coverage for business logic
5. **Use Structured Logging**: Include relevant context in all log messages

Example handler function:

```go
// GetUsers retrieves a paginated list of users
func GetUsers(c *fiber.Ctx) error {
    users, total, err := userService.GetUsers(c.Context())
    if err != nil {
        return response.Error(c).
            Message("Failed to retrieve users").
            Code("USER_FETCH_ERROR").
            Details(err.Error()).
            Send()
    }

    return response.Success(c).
        Message("Users retrieved successfully").
        Data(users).
        Pagination(1, 10, total).
        Send()
    }
```

### Code Quality

- **Go**: Use `go fmt`, `go vet`, `golangci-lint` and `gosec`:)
- **TypeScript**: Use Prettier and ESLint for code formatting
- **Testing**: Maintain test coverage above 80%
- **Documentation**: Update README files for package changes

## License

This project is licensed under the terms specified in the LICENSE file.
