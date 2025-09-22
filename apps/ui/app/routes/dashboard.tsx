import { Dashboard } from '~/components/dashboard/dashboard';
import { useCurrentUser } from '~/hooks';
import { Navigate } from 'react-router';

export function meta() {
  return [
    { title: 'Dashboard | PWS ELO' },
    {
      name: 'description',
      content: 'Je persoonlijke dashboard met al je vakken',
    },
  ];
}

export default function DashboardPage() {
  const { data: user, isLoading } = useCurrentUser();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Dashboard />;
}
