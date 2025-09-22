import { Dashboard } from '~/components/dashboard/dashboard';
import { ProtectedRoute } from '~/components';

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
  return (
    <ProtectedRoute>
      <Dashboard />
    </ProtectedRoute>
  );
}
