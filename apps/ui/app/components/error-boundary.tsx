import React from 'react';
import { AlertTriangle, RefreshCw, Home } from 'lucide-react';
import { Button } from './ui/button';

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

interface ErrorBoundaryProps {
  children: React.ReactNode;
  fallback?: React.ComponentType<{ error: Error; resetError: () => void }>;
}

export class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error boundary caught an error:', error, errorInfo);

    // You can also log the error to an error reporting service here
    // Example: Sentry.captureException(error, { extra: errorInfo });
  }

  resetError = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        const FallbackComponent = this.props.fallback;
        return <FallbackComponent error={this.state.error!} resetError={this.resetError} />;
      }

      return <DefaultErrorFallback error={this.state.error!} resetError={this.resetError} />;
    }

    return this.props.children;
  }
}

interface ErrorFallbackProps {
  error: Error;
  resetError: () => void;
}

function DefaultErrorFallback({ error, resetError }: ErrorFallbackProps) {
  const isDevelopment = process.env.NODE_ENV === 'development';

  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 px-4">
      <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-6 text-center">
        <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-error-100 mb-4">
          <AlertTriangle className="h-6 w-6 text-error-600" />
        </div>

        <h1 className="text-xl font-bold text-neutral-900 mb-2">Er is iets misgegaan</h1>

        <p className="text-neutral-600 mb-6">
          Er is een onverwachte fout opgetreden. Probeer de pagina te vernieuwen of ga terug naar de
          startpagina.
        </p>

        {isDevelopment && (
          <div className="mb-6 p-4 bg-neutral-100 rounded-lg text-left">
            <p className="text-sm font-medium text-neutral-700 mb-2">
              Foutdetails (alleen zichtbaar in ontwikkelingsmodus):
            </p>
            <p className="text-xs text-neutral-600 font-mono break-all">{error.message}</p>
            {error.stack && (
              <pre className="text-xs text-neutral-600 mt-2 overflow-auto max-h-32">
                {error.stack}
              </pre>
            )}
          </div>
        )}

        <div className="flex flex-col sm:flex-row gap-3">
          <Button onClick={resetError} className="flex-1" variant="outline">
            <RefreshCw className="h-4 w-4 mr-2" />
            Probeer opnieuw
          </Button>

          <Button onClick={() => (window.location.href = '/')} className="flex-1">
            <Home className="h-4 w-4 mr-2" />
            Naar startpagina
          </Button>
        </div>
      </div>
    </div>
  );
}

// Hook for functional components to handle errors
export function useErrorHandler() {
  return (error: Error, errorInfo?: React.ErrorInfo) => {
    console.error('Error caught by error handler:', error, errorInfo);

    // You can also report to error tracking service here
    // Example: Sentry.captureException(error, { extra: errorInfo });
  };
}

// HOC to wrap components with error boundary
export function withErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  fallback?: React.ComponentType<{ error: Error; resetError: () => void }>
) {
  const WrappedComponent = (props: P) => (
    <ErrorBoundary fallback={fallback}>
      <Component {...props} />
    </ErrorBoundary>
  );

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;

  return WrappedComponent;
}
