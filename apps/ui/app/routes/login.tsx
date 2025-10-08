import { LoginForm } from '~/components/auth/login-form';
import { FeatureDisabled } from '~/components/feature-disabled';
import { useAuth } from '~/hooks';
import { useNavigate } from 'react-router';
import { useEffect } from 'react';
import { env } from '~/lib/env';

export function meta() {
  return [
    { title: 'Inloggen | PWS ELO' },
    { name: 'description', content: 'Log in op je PWS ELO account' },
  ];
}

export default function Login() {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  // Handle navigation when authenticated
  useEffect(() => {
    if (isAuthenticated && !isLoading) {
      navigate('/dashboard', { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

  // Don't render login form if already authenticated (navigation will happen via useEffect)
  if (isAuthenticated && !isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-neutral-50">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  // Check if login feature is disabled
  if (!env.features.enableLogin) {
    return (
      <FeatureDisabled
        featureName="Inloggen"
        description="De inlog functionaliteit is momenteel uitgeschakeld door de beheerder."
      />
    );
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
        <LoginForm isLoading={isLoading} />
      </div>
    </div>
  );
}
