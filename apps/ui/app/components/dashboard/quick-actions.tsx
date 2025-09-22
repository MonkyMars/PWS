import { Bell, FileText, Calendar, ExternalLink } from 'lucide-react';
import { Button } from '~/components/ui/button';

export function QuickActions() {
  return (
    <div className="bg-white rounded-lg border border-neutral-200 p-6">
      <h3 className="text-lg font-semibold text-neutral-900 mb-4">Snelle Acties</h3>

      <div className="space-y-3">
        <Button variant="outline" className="w-full justify-start" size="sm">
          <Bell className="h-4 w-4 mr-2" />
          Alle Mededelingen
        </Button>

        <Button variant="outline" className="w-full justify-start" size="sm">
          <FileText className="h-4 w-4 mr-2" />
          Recente Bestanden
        </Button>

        <Button variant="outline" className="w-full justify-start" size="sm">
          <Calendar className="h-4 w-4 mr-2" />
          Rooster Bekijken
        </Button>

        <a
          href="#"
          className="flex items-center justify-start w-full px-3 py-2 text-sm font-medium text-neutral-700 border border-neutral-300 rounded-lg hover:bg-neutral-50 transition-colors"
        >
          <ExternalLink className="h-4 w-4 mr-2" />
          Help & Ondersteuning
        </a>
      </div>
    </div>
  );
}
