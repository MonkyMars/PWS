# Components

This directory contains all reusable UI components for the PWS ELO application.

## Directory Structure

```
components/
├── auth/               # Authentication related components
├── dashboard/          # Dashboard specific components
├── files/              # File handling and viewing components
├── subjects/           # Subject management components
├── ui/                 # Basic UI primitives (buttons, inputs, etc.)
├── footer.tsx          # Application footer
└── navigation.tsx      # Main navigation component
```

## Component Categories

### Authentication (`auth/`)

- `login-form.tsx` - User login form with validation
- `register-form.tsx` - User registration form with Zod validation
- `validation.ts` - Form validation schemas

### Dashboard (`dashboard/`)

- `dashboard.tsx` - Main dashboard layout
- `subject-card.tsx` - Subject card component for dashboard
- `quick-actions.tsx` - Quick action buttons sidebar
- `recent-activity.tsx` - Recent activity feed

### Files (`files/`)

- `file-viewer.tsx` - In-app file viewer with zoom and download options

### Subjects (`subjects/`)

- `subject-detail.tsx` - Detailed subject view with announcements and files

### UI Primitives (`ui/`)

- `button.tsx` - Customizable button component
- `input.tsx` - Form input with label and error handling

## Design Principles

- **Accessibility**: All components follow WCAG guidelines
- **Responsiveness**: Mobile-first responsive design
- **Type Safety**: Full TypeScript coverage with proper prop types
- **Consistency**: Unified design system across all components
- **Performance**: Optimized rendering with React best practices

## Usage Guidelines

1. Always use TypeScript interfaces for component props
2. Include proper accessibility attributes (ARIA, etc.)
3. Follow the established design tokens from `app.css`
4. Use the custom hooks from `../hooks` for data fetching
5. Implement proper error boundaries and loading states
