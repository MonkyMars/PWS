import {
  Bell,
  FileText,
  Calendar,
  ExternalLink,
  LayoutGridIcon,
  LucideChartColumnStacked,
} from 'lucide-react';
import { Button } from '~/components/ui/button';
import { useAuth } from '~/hooks';

// TODO: Add real quick links and functionality
const StudentQuickActions = () => {
  return (
    <div className="bg-white rounded-lg border border-neutral-200 p-6">
      <h3 className="text-lg font-semibold text-neutral-900 mb-4">Snelle Acties</h3>

      <div className="space-y-3">
        <Button variant="outline" className="w-full justify-start" size="sm">
          <LucideChartColumnStacked className="h-4 w-4 mr-2" />
          Mijn Opdrachten
        </Button>
      </div>
    </div>
  );
};

// TODO: Add real quick links and functionality
const TeacherQuickActions = () => {
  return (
    <div className="bg-white rounded-lg border border-neutral-200 p-6">
      <h3 className="text-lg font-semibold text-neutral-900 mb-4">Snelle Acties</h3>

      <div className="space-y-3">
        <Button variant="outline" className="w-full justify-start" size="sm">
          <LayoutGridIcon className="h-4 w-4 mr-2" />
          Klassen Overzicht
        </Button>

        <Button variant="outline" className="w-full justify-start" size="sm">
          <FileText className="h-4 w-4 mr-2" />
          Nieuw Bestand Uploaden
        </Button>

        <Button variant="outline" className="w-full justify-start" size="sm">
          <Calendar className="h-4 w-4 mr-2" />
          Ingeleverde Opdrachten
        </Button>
      </div>
    </div>
  );
};

export function QuickActions() {
  const { user } = useAuth();
  if (!user) return null;
  if (user.role === 'teacher' || user.role === 'admin') {
    return <TeacherQuickActions />;
  }
  return <StudentQuickActions />;
}
