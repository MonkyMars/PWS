import { LoginForm } from '~/components/auth/login-form';
import { useAuth } from '~/hooks';
import { Navigate } from 'react-router';

export function meta() {
  return [
    { title: 'Inloggen | PWS ELO' },
    { name: 'description', content: 'Log in op je PWS ELO account' },
  ];
}

export default function Login() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-bold text-neutral-900">
            Inloggen op je account
          </h2>
          <p className="mt-2 text-center text-sm text-neutral-600">
            Of{' '}
            <a href="/register" className="font-medium text-primary-600 hover:text-primary-500">
              maak een nieuw account aan
            </a>
          </p>
        </div>
        <LoginForm />
      </div>
    </div>
  );
}
