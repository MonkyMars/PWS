# PWS UI Application

The user interface for the PWS (Personal Web Space) application, built with React Router v7, TypeScript, and Tailwind CSS. This application provides a modern, responsive web interface for interacting with the PWS API.

## Overview

This React application features:

- Modern React with TypeScript for type safety
- React Router v7 for routing and navigation
- Tailwind CSS for styling and responsive design
- Server-side rendering (SSR) capabilities
- Dark mode support
- Responsive design optimized for all device sizes

## Technology Stack

- **React 19** - UI library with latest features
- **React Router v7** - File-based routing with SSR support
- **TypeScript** - Type safety and better developer experience
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
import { type RouteConfig, index } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"), // Maps "/" to home.tsx
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
export default [index("routes/home.tsx")] satisfies RouteConfig;
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
    throw new Error(result.error?.message || "Failed to fetch user");
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
