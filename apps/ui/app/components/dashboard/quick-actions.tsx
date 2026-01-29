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

// Student quick actions with real navigation
import { useNavigate } from 'react-router';

const StudentQuickActions = () => {
  const navigate = useNavigate();
  return (
    <div className="bg-white rounded-lg border border-neutral-200 p-6">
      <h3 className="text-lg font-semibold text-neutral-900 mb-4">Snelle Acties</h3>

      <div className="space-y-3">
        <Button
          variant="outline"
          className="w-full justify-start"
          size="sm"
          onClick={() => navigate('/deadlines/mine')}
        >
          <LucideChartColumnStacked className="h-4 w-4 mr-2" />
          Mijn Opdrachten
        </Button>
      </div>
    </div>
  );
};

// Teacher quick actions with real navigation
const TeacherQuickActions = () => {
  const navigate = useNavigate();
  return (
    <div className="bg-white rounded-lg border border-neutral-200 p-6">
      <h3 className="text-lg font-semibold text-neutral-900 mb-4">Snelle Acties</h3>

      <div className="space-y-3">
        <Button
          variant="outline"
          className="w-full justify-start"
          size="sm"
          onClick={() => navigate('/classes')}
        >
          <LayoutGridIcon className="h-4 w-4 mr-2" />
          Klassen Overzicht
        </Button>

        <Button
          variant="outline"
          className="w-full justify-start"
          size="sm"
          onClick={() => navigate('/files/upload')}
        >
          <FileText className="h-4 w-4 mr-2" />
          Nieuw Bestand Uploaden
        </Button>

        <Button
          variant="outline"
          className="w-full justify-start"
          size="sm"
          onClick={() => navigate('/deadlines/submissions')}
        >
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
