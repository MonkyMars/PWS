# App Root

This directory contains the main application code for the PWS ELO frontend.

## Directory Structure

```
app/
├── components/          # Reusable UI components
├── hooks/              # Custom React hooks for data fetching and state management
├── lib/                # Utility libraries and configurations
├── routes/             # React Router v7 route components
├── types/              # TypeScript type definitions
├── app.css             # Global styles and theme variables
└── root.tsx            # Root application component with providers
```

## Key Features

- **TypeScript**: Fully typed codebase with strict type checking
- **React Router v7**: Modern routing with file-based routing conventions
- **TanStack Query**: Efficient data fetching and caching
- **Tailwind CSS**: Utility-first CSS framework with custom theme
- **Zod**: Runtime type validation for forms and API responses

## Getting Started

1. Install dependencies: `bun install`
2. Start development server: `bun dev`
3. Build for production: `bun build`

## Architecture Principles

- **Component-based**: Modular, reusable components
- **Type-safe**: Comprehensive TypeScript coverage
- **Performance-focused**: Optimized loading and caching strategies
- **Accessible**: WCAG compliant UI components
- **Mobile-first**: Responsive design for all devices
