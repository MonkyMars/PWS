import { type ReactNode, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { useAuth } from '~/hooks/use-auth-context';

interface ProtectedRouteProps {
  children: ReactNode;
  redirectTo?: string;
  requiredRole?: string;
}

export function ProtectedRoute({
  children,
  redirectTo = '/login',
  requiredRole,
}: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  // Handle navigation when authentication state changes
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate(redirectTo, { state: { from: location }, replace: true });
    }
  }, [isAuthenticated, isLoading, navigate, redirectTo, location]);

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  // Don't render children if not authenticated (navigation will happen via useEffect)
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  // Check role-based access if required
  if (requiredRole && user?.role !== requiredRole) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Access Denied</h1>
          <p className="text-gray-600 mb-6">You don't have permission to access this page.</p>
          <button
            onClick={() => window.history.back()}
            className="text-primary-600 hover:text-primary-500"
          >
            Go back
          </button>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
