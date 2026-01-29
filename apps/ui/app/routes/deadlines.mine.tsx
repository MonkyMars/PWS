import React, { useMemo, useState } from 'react';
import { useDeadlines } from '~/hooks/use-deadlines';
import { useAuth } from '~/hooks';
import { Input } from '~/components/ui/input';
import { Button } from '~/components/ui/button';
import { FileText, AlertCircle, Search, UploadCloud } from 'lucide-react';
import { format } from 'date-fns';
import { useNavigate } from 'react-router';
import { useStudentSubmission } from '~/hooks/use-student-submission';

export default function StudentDeadlinesPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { data: deadlines = [], isLoading } = useDeadlines();
  const [search, setSearch] = useState('');

  // Filter deadlines by title, subject, or code
  const filtered = useMemo(() => {
    if (!search.trim()) return deadlines;
    const term = search.toLowerCase();
    return deadlines.filter((d) => {
      return (
        d.title.toLowerCase().includes(term) ||
        d.subject?.name?.toLowerCase().includes(term) ||
        d.subject?.code?.toLowerCase().includes(term)
      );
    });
  }, [search, deadlines]);

  if (!user || user.role !== 'student') {
    return (
      <div className="p-8 text-center text-red-600">
        <AlertCircle className="inline-block mr-2" />
        Je hebt geen toegang tot deze pagina.
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto py-10 px-4">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">Mijn Opdrachten</h1>
        <Button variant="outline" onClick={() => navigate('/dashboard')}>
          Terug naar dashboard
        </Button>
      </div>

      <div className="mb-6 flex items-center gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-neutral-400" />
          <Input
            type="text"
            placeholder="Zoek op opdracht, vak, code..."
            className="w-full border border-neutral-300 rounded-lg pl-10 pr-4 py-2"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      {isLoading ? (
        <div className="text-neutral-500">Laden...</div>
      ) : filtered.length === 0 ? (
        <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
          <FileText className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-neutral-900 mb-2">Geen opdrachten gevonden</h3>
          <p className="text-neutral-600">
            Er zijn geen opdrachten gevonden die overeenkomen met je zoekopdracht.
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full bg-white rounded-lg border border-neutral-200">
            <thead>
              <tr>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Opdracht
                </th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">Vak</th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Deadline
                </th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Status
                </th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Acties
                </th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((d) => {
                // Fetch submission for each deadline
                // eslint-disable-next-line react-hooks/rules-of-hooks
                const { data: submission, isLoading: subLoading } = useStudentSubmission(d.id);

                let statusNode: React.ReactNode = (
                  <span className="text-xs px-2 py-1 bg-neutral-100 text-neutral-700 rounded">
                    Nog niet ingeleverd
                  </span>
                );
                if (subLoading) {
                  statusNode = (
                    <span className="text-xs px-2 py-1 bg-neutral-100 text-neutral-700 rounded">
                      Laden...
                    </span>
                  );
                } else if (submission) {
                  if (submission.isLate) {
                    statusNode = (
                      <span className="text-xs px-2 py-1 bg-red-100 text-red-700 rounded">
                        Te laat
                      </span>
                    );
                  } else {
                    statusNode = (
                      <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded">
                        Op tijd ingeleverd
                      </span>
                    );
                  }
                  if (submission.isUpdated) {
                    statusNode = (
                      <>
                        {statusNode}
                        <span className="ml-2 text-xs px-2 py-1 bg-yellow-100 text-yellow-800 rounded">
                          Gewijzigd na deadline
                        </span>
                      </>
                    );
                  }
                }

                return (
                  <tr key={d.id} className="border-t border-neutral-100 hover:bg-neutral-50">
                    <td className="px-4 py-2">{d.title}</td>
                    <td className="px-4 py-2">
                      {d.subject?.name}{' '}
                      <span className="text-xs text-neutral-400">({d.subject?.code})</span>
                    </td>
                    <td className="px-4 py-2">{format(new Date(d.dueDate), 'dd-MM-yyyy HH:mm')}</td>
                    <td className="px-4 py-2">{statusNode}</td>
                    <td className="px-4 py-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => navigate(`/deadlines/${d.id}/submission`)}
                        title="Lever in / bekijk opdracht"
                      >
                        <UploadCloud className="h-4 w-4 mr-1" />
                        {submission ? 'Bekijk / Wijzig' : 'Inleveren'}
                      </Button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
