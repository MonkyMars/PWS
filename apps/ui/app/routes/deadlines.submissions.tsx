import { useMemo, useState } from 'react';
import { Input } from '~/components/ui/input';
import { Button } from '~/components/ui/button';
import { Printer, Search, FileText, AlertCircle } from 'lucide-react';
import { useAuth } from '~/hooks';
import { useNavigate } from 'react-router';
import { format } from 'date-fns';
import { useTeacherSubmissions } from '~/hooks/use-teacher-submissions';

// Print utility (replace with real print logic)
function printSubmission(submission: any) {
  // For demo: just alert, replace with real print logic
  window.print();
}

export default function TeacherSubmissionsPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { data: submissions = [], isLoading: loading } = useTeacherSubmissions();
  const [search, setSearch] = useState('');

  // Filter submissions by student name/email or deadline title/subject
  const filtered = useMemo(() => {
    if (!search.trim()) return submissions;
    const term = search.toLowerCase();
    return submissions.filter((s) => {
      return (
        s.student?.name.toLowerCase().includes(term) ||
        s.student?.email.toLowerCase().includes(term) ||
        s.deadline?.title.toLowerCase().includes(term) ||
        s.deadline?.subject.name.toLowerCase().includes(term) ||
        s.deadline?.subject.code.toLowerCase().includes(term)
      );
    });
  }, [search, submissions]);

  if (!user || (user.role !== 'teacher' && user.role !== 'admin')) {
    return (
      <div className="p-8 text-center text-red-600">
        <AlertCircle className="inline-block mr-2" />
        Je hebt geen toegang tot deze pagina.
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto py-10 px-4">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold">Ingeleverde Opdrachten</h1>
        <Button variant="outline" onClick={() => window.print()}>
          <Printer className="h-4 w-4 mr-2" />
          Print alles
        </Button>
      </div>

      <div className="mb-6 flex items-center gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-neutral-400" />
          <Input
            type="text"
            placeholder="Zoek op student, opdracht, vak..."
            className="w-full border border-neutral-300 rounded-lg pl-10 pr-4 py-2"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      {loading ? (
        <div className="text-neutral-500">Laden...</div>
      ) : filtered.length === 0 ? (
        <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
          <FileText className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-neutral-900 mb-2">
            Geen ingeleverde opdrachten gevonden
          </h3>
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
                  Student
                </th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Opdracht
                </th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">Vak</th>
                <th className="px-4 py-2 text-left text-sm font-semibold text-neutral-700">
                  Ingeleverd op
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
              {filtered.map((s) => (
                <tr key={s.id} className="border-t border-neutral-100 hover:bg-neutral-50">
                  <td className="px-4 py-2">
                    <div className="font-medium">{s.student?.name}</div>
                    <div className="text-xs text-neutral-500">{s.student?.email}</div>
                  </td>
                  <td className="px-4 py-2">
                    <div>{s.deadline?.title}</div>
                  </td>
                  <td className="px-4 py-2">
                    <div>
                      {s.deadline?.subject.name}{' '}
                      <span className="text-xs text-neutral-400">({s.deadline?.subject.code})</span>
                    </div>
                  </td>
                  <td className="px-4 py-2">{format(new Date(s.createdAt), 'dd-MM-yyyy HH:mm')}</td>
                  <td className="px-4 py-2">
                    {s.isLate ? (
                      <span className="text-xs px-2 py-1 bg-red-100 text-red-700 rounded">
                        Te laat
                      </span>
                    ) : (
                      <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded">
                        Op tijd
                      </span>
                    )}
                    {s.isUpdated && (
                      <span className="ml-2 text-xs px-2 py-1 bg-yellow-100 text-yellow-800 rounded">
                        Gewijzigd na deadline
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-2 flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => printSubmission(s)}
                      title="Print deze opdracht"
                    >
                      <Printer className="h-4 w-4" />
                    </Button>
                    {/* Optionally: link to view/download files */}
                    {s.fileIds && s.fileIds.length > 0 && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => {
                          // Replace with real file download/view logic
                          alert('Bestanden openen: ' + s.fileIds.join(', '));
                        }}
                        title="Bekijk bestanden"
                      >
                        <FileText className="h-4 w-4" />
                      </Button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
