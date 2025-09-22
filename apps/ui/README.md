# PWS ELO Frontend

A modern, TypeScript-based frontend for the PWS Electronic Learning Environment (ELO) built with React Router v7, TanStack Query, and Tailwind CSS.

## 🚀 Quick Start

```bash
# Install dependencies
bun install

# Start development server
bun dev

# Build for production
bun build

# Type checking
bun typecheck
```

## 📁 Project Structure

```
apps/ui/
├── app/                    # Main application code
│   ├── components/         # Reusable UI components
│   │   ├── auth/          # Authentication forms
│   │   ├── dashboard/     # Dashboard components
│   │   ├── files/         # File viewer and management
│   │   ├── subjects/      # Subject detail views
│   │   └── ui/            # Basic UI primitives
│   ├── hooks/             # Custom React hooks
│   ├── lib/               # Utilities and configurations
│   ├── routes/            # Page components
│   ├── types/             # TypeScript type definitions
│   ├── app.css            # Global styles and theme
│   └── root.tsx           # App root with providers
├── public/                # Static assets
├── package.json           # Dependencies and scripts
├── tsconfig.json          # TypeScript configuration
└── vite.config.ts         # Vite build configuration
```

## 🛠 Technology Stack

- **Framework**: React 19 with React Router v7
- **Language**: TypeScript with strict type checking
- **Styling**: Tailwind CSS with custom design system
- **Data Fetching**: TanStack Query for caching and synchronization
- **Validation**: Zod for runtime type validation
- **Build Tool**: Vite for fast development and builds
- **Package Manager**: Bun for fast dependency management

## 🎨 Design System

### Colors

- **Primary**: Educational blue theme (`#3b82f6`)
- **Secondary**: Warm accent orange (`#f27318`)
- **Success**: Green (`#22c55e`)
- **Warning**: Yellow (`#f59e0b`)
- **Error**: Red (`#ef4444`)

### Typography

- **Font Family**: Inter (Google Fonts)
- **Scale**: Tailwind's default type scale
- **Weights**: 400 (normal), 500 (medium), 600 (semibold), 700 (bold)

### Components

All components follow accessibility guidelines and include:

- Proper ARIA attributes
- Keyboard navigation support
- Focus management
- Screen reader compatibility

## 🔒 Authentication

The app uses token-based authentication with:

- JWT tokens stored in localStorage
- Automatic token refresh
- Protected route guards
- Role-based access control (student, teacher, admin)

## 📱 Features

### For Students

- **Dashboard**: Overview of all enrolled subjects
- **Subject Details**: Announcements, files, and resources
- **File Viewer**: In-app preview of documents and images
- **Mobile-First**: Optimized for all device sizes

### For Teachers (Future)

- File upload and management
- Announcement creation
- Student progress tracking

### For Admins (Future)

- User management
- System configuration
- Analytics and reporting

## 🌐 Internationalization

The application is currently built for Dutch users with:

- Dutch language interface
- Dutch date/time formatting
- Netherlands-specific form validation
- Local cultural conventions

## 🔧 Development

### Code Quality

- **ESLint**: Configured for React and TypeScript
- **TypeScript**: Strict mode with no `any` types allowed
- **Prettier**: Consistent code formatting
- **Git Hooks**: Pre-commit linting and type checking

### Testing Strategy

- Unit tests for utility functions
- Component testing with React Testing Library
- Integration tests for critical user flows
- E2E tests for complete user journeys

### Performance

- **Code Splitting**: Route-based lazy loading
- **Image Optimization**: Automatic image optimization
- **Bundle Analysis**: Regular bundle size monitoring
- **Caching**: Intelligent data caching with TanStack Query

## 🚀 Deployment

### Environment Variables

```bash
API_URL=http://localhost:8080/api  # Backend API URL
NODE_ENV=production                # Environment mode
```

### Build Output

- Static files in `build/` directory
- Optimized for CDN deployment
- Service worker for offline support
- Progressive Web App capabilities

## 🤝 Contributing

1. Follow the established TypeScript patterns
2. Use the custom hooks for data fetching
3. Implement proper error boundaries
4. Include accessibility attributes
5. Write meaningful commit messages
6. Update documentation for new features

## 📚 Learning Resources

- [React Router v7 Documentation](https://reactrouter.com/dev)
- [TanStack Query Guide](https://tanstack.com/query/latest)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs)

---

Built with ❤️ for PWS School students and teachers.

- **Tailwind CSS v4** - Utility-first CSS framework
- **Vite** - Fast build tool and development server
- **Bun** - Package manager and runtime

## Project Structure

```
apps/ui/
├── app/                    # Application source code
│   ├── routes/            # Route components
│   │   └── home.tsx       # Home page route
│   ├── welcome/           # Welcome component and assets
│   │   ├── welcome.tsx    # Welcome page component
│   │   ├── logo-dark.svg  # Dark theme logo
│   │   └── logo-light.svg # Light theme logo
│   ├── app.css           # Global styles
│   ├── root.tsx          # Root layout component
│   └── routes.ts         # Route configuration
├── public/               # Static assets
│   └── favicon.ico       # Application favicon
├── build/               # Build output (generated)
├── package.json         # Dependencies and scripts
├── vite.config.ts       # Vite configuration
├── tsconfig.json        # TypeScript configuration
├── react-router.config.ts # React Router configuration
├── Dockerfile           # Container configuration
├── bun.lock            # Lock file for dependencies
└── README.md           # This documentation
```

## Getting Started

### Prerequisites

- [Bun](https://bun.sh/) (recommended) or Node.js 18+
- Modern web browser with JavaScript enabled

### Installation

1. Navigate to the UI directory:

   ```bash
   cd apps/ui
   ```

2. Install dependencies:
   ```bash
   bun install
   ```

### Development

Start the development server:

```bash
bun run dev
```

The application will be available at `http://localhost:5173` with hot module replacement enabled.

### Building

Build the application for production:

```bash
bun run build
```

The built application will be in the `build/` directory.

### Production Server

Start the production server:

```bash
bun run start
```

This runs the built application with server-side rendering.

### Type Checking

Run TypeScript type checking:

```bash
bun run typecheck
```

## Application Architecture

### Routing

The application uses React Router v7 with file-based routing:

```typescript
// app/routes.ts
import { type RouteConfig, index } from '@react-router/dev/routes';

export default [
  index('routes/home.tsx'), // Maps "/" to home.tsx
] satisfies RouteConfig;
```

### Components

#### Root Layout (`app/root.tsx`)

The root layout component provides:

- HTML document structure
- Meta tags and external links
- Font loading optimization
- Global styles and scripts

#### Home Route (`app/routes/home.tsx`)

The home page route that renders the Welcome component.

#### Welcome Component (`app/welcome/welcome.tsx`)

The main welcome screen featuring:

- Responsive logo display with dark/light mode support
- Navigation links to external resources
- Tailwind CSS styling with hover effects
- Accessibility features

### Styling

The application uses Tailwind CSS for styling with:

- **Responsive Design**: Mobile-first responsive breakpoints
- **Dark Mode**: Automatic dark mode support with `dark:` classes
- **Component Classes**: Utility classes for rapid development
- **Custom Fonts**: Inter font family for better typography

Example usage:

```tsx
<div className="flex items-center justify-center pt-16 pb-4">
  <img src={logoLight} className="block w-full dark:hidden" alt="PWS Logo" />
</div>
```

### TypeScript Integration

The application uses TypeScript throughout:

```typescript
// Type-safe component props
interface ResourceLink {
  href: string;
  text: string;
  icon: React.ReactElement;
}

// Type-safe route configuration
export default [index('routes/home.tsx')] satisfies RouteConfig;
```

## Development Guidelines

### Component Development

1. **Functional Components**: Use functional components with hooks
2. **TypeScript**: Always type component props and state
3. **JSDoc Comments**: Document component purpose and props
4. **Accessibility**: Include proper ARIA labels and semantic HTML

Example component structure:

```tsx
/**
 * Component description and purpose.
 *
 * @param props - Component properties
 * @returns JSX element
 */
export function MyComponent({ title, onClick }: MyComponentProps) {
  return (
    <button
      onClick={onClick}
      className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      aria-label={`Button: ${title}`}
    >
      {title}
    </button>
  );
}

interface MyComponentProps {
  title: string;
  onClick: () => void;
}
```

### Styling Guidelines

1. **Utility Classes**: Prefer Tailwind utility classes over custom CSS
2. **Responsive Design**: Always consider mobile-first design
3. **Dark Mode**: Include dark mode variants where appropriate
4. **Consistent Spacing**: Use Tailwind's spacing scale

### Performance Considerations

1. **Code Splitting**: Leverage React Router's built-in code splitting
2. **Image Optimization**: Use appropriate image formats and sizes
3. **Font Loading**: Optimize font loading with preconnect links
4. **Bundle Analysis**: Monitor bundle size and dependencies

## API Integration

When integrating with the PWS API:

1. **Type Safety**: Define TypeScript interfaces for API responses
2. **Error Handling**: Implement proper error boundaries and user feedback
3. **Loading States**: Show appropriate loading indicators
4. **Caching**: Consider implementing client-side caching strategies

Example API integration:

```typescript
interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

async function fetchUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  const result: ApiResponse<User> = await response.json();

  if (!result.success) {
    throw new Error(result.error?.message || 'Failed to fetch user');
  }

  return result.data!;
}
```

## Deployment

### Docker Deployment

The application includes a Dockerfile for containerized deployment:

```bash
# Build the container
docker build -t pws-ui .

# Run the container
docker run -p 3000:3000 pws-ui
```

### Environment Variables

Configure the application using environment variables:

- `VITE_API_URL` - Backend API URL
- `VITE_APP_TITLE` - Application title
- `NODE_ENV` - Environment (development/production)

### Platform Deployment

The containerized application can be deployed to any platform that supports Docker:

- AWS ECS
- Google Cloud Run
- Azure Container Apps
- Digital Ocean App Platform
- Fly.io
- Railway

## Browser Support

The application supports:

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Contributing

1. Follow the existing code style and patterns
2. Write TypeScript for all new code
3. Include JSDoc comments for components and functions
4. Test responsive design on multiple screen sizes
5. Ensure accessibility compliance

## Dependencies

### Core Dependencies

- `react` - UI library
- `react-dom` - DOM rendering
- `react-router` - Routing and navigation
- `@react-router/node` - Node.js adapter
- `@react-router/serve` - Production server

### Development Dependencies

- `@react-router/dev` - Development tools
- `typescript` - Type checking
- `vite` - Build tool
- `tailwindcss` - CSS framework
- `@tailwindcss/vite` - Vite integration

## License

This project is part of the PWS application suite.
