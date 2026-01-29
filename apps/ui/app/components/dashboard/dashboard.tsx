import {
  BookOpen,
  X,
  Search,
  CornerDownLeft,
  Command,
  Check,
  CalendarClockIcon,
  ArrowRight,
} from 'lucide-react';
import { SubjectCard } from './subject-card';
import { QuickActions } from './quick-actions';
import { useCurrentUser, useSubjects, useDebounce } from '~/hooks';
import { Input } from '../ui/input';
import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router';
import Restricted from '../restricted';
import { Button } from '../ui/button';
import { useDeadlines } from '~/hooks/use-deadlines';
import { Link } from 'react-router';

export function Dashboard() {
  const { data: user } = useCurrentUser();
  const { data: subjects, isLoading: subjectsLoading } = useSubjects();
  const { data: deadlines, isLoading: deadlinesLoading } = useDeadlines();
  const [subjectSearchValue, setSubjectSearchValue] = useState<string>('');
  const [deadlineSearchValue, setDeadlineSearchValue] = useState<string>('');
  const navigate = useNavigate();
  const searchInputRef = useRef<HTMLInputElement>(null);
  const [managingSubjects, setManagingSubjects] = useState<boolean>(false);

  // Debounce search value for better performance
  const debouncedSubjectSearchValue = useDebounce(subjectSearchValue, 200);
  const debouncedDeadlineSearchValue = useDebounce(deadlineSearchValue, 200);

  // Filter and sort subjects based on debounced search with prioritized scoring
  const filteredSubjects = (() => {
    if (!subjects) return [];
    if (!debouncedSubjectSearchValue.trim()) return subjects;

    const searchTerm = debouncedSubjectSearchValue.toLowerCase();

    // Score subjects based on match quality
    const scoredSubjects = subjects
      .map((subject) => {
        const name = subject.name.toLowerCase();
        const code = subject.code.toLowerCase();

        let score = 0;

        // Name field (highest priority)
        if (name.startsWith(searchTerm)) {
          score += 100;
        } else if (name.includes(searchTerm)) {
          score += 10;
        }

        // Code field (medium priority)
        if (code.startsWith(searchTerm)) {
          score += 50;
        } else if (code.includes(searchTerm)) {
          score += 5;
        }

        return { subject, score };
      })
      .filter(({ score }) => score > 0)
      .sort((a, b) => b.score - a.score)
      .map(({ subject }) => subject);

    return scoredSubjects;
  })();

  const filteredDeadlines = (() => {
    if (!deadlines) return [];
    if (!debouncedDeadlineSearchValue.trim()) return deadlines;

    const searchTerm = debouncedDeadlineSearchValue.toLowerCase();

    // Score subjects based on match quality
    const scoredDeadlines = deadlines
      .map((deadline) => {
        const name = deadline.title.toLowerCase();
        const code = deadline.subject.code.toLowerCase();
        const subject = deadline.subject.name.toLowerCase();

        let score = 0;

        // Name field (highest priority)
        if (name.startsWith(searchTerm)) {
          score += 100;
        } else if (name.includes(searchTerm)) {
          score += 10;
        }

        // Subject field (medium priority)
        if (subject.startsWith(searchTerm)) {
          score += 50;
        } else if (subject.includes(searchTerm)) {
          score += 5;
        }

        // Code field (low priority)
        if (code.startsWith(searchTerm)) {
          score += 25;
        } else if (code.includes(searchTerm)) {
          score += 3;
        }

        return { deadline, score };
      })
      .filter(({ score }) => score > 0)
      .sort((a, b) => b.score - a.score)
      .map(({ deadline }) => deadline);

    return scoredDeadlines;
  })();

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Goedemorgen';
    if (hour < 17) return 'Goedemiddag';
    return 'Goedenavond';
  };

  const handleSearchKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Escape') {
      setSubjectSearchValue('');
      searchInputRef.current?.blur();
    } else if (e.key === 'Enter' && filteredSubjects.length > 0) {
      // Navigate to the first subject in the filtered results
      navigate(`/subjects/${filteredSubjects[0].id}`);
    }
  };

  // Global keyboard shortcut to focus search and navigate to subjects
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      const target = e.target as HTMLElement;
      const isInInput = ['INPUT', 'TEXTAREA'].includes(target?.tagName);

      // Focus search with 'S'
      if (e.key.toLowerCase() === 's' && !isInInput && searchInputRef.current) {
        e.preventDefault();
        searchInputRef.current.focus();
        return;
      }

      // Navigate to subjects with number keys (1-9)
      const keyNumber = parseInt(e.key);
      if (
        !isNaN(keyNumber) &&
        keyNumber >= 0 &&
        keyNumber <= 9 &&
        !isInInput &&
        filteredSubjects.length > 0
      ) {
        const subjectIndex = keyNumber;
        if (filteredSubjects[subjectIndex]) {
          e.preventDefault();
          navigate(`/subjects/${filteredSubjects[subjectIndex].id}`);
        }
      }
    };

    document.addEventListener('keydown', handleGlobalKeyDown);
    return () => document.removeEventListener('keydown', handleGlobalKeyDown);
  }, [filteredSubjects, navigate]);

  return (
    <div className="min-h-screen bg-neutral-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Welcome Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-neutral-900 mb-2">
            {getGreeting()}, {user?.username}!
          </h1>
          <p className="text-neutral-600">
            Welkom terug in je PWS ELO dashboard. Hier vind je al je vakken en recente activiteiten.
          </p>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-lg p-6 border border-neutral-200">
            <div className="flex items-center">
              <div className="p-2 bg-primary-100 rounded-lg">
                <BookOpen className="h-6 w-6 text-primary-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-neutral-600">Vakken</p>
                <p className="text-2xl font-bold text-neutral-900">{subjects?.length || 0}</p>
              </div>
            </div>
          </div>
          <div className="bg-white rounded-lg p-6 border border-neutral-200">
            <div className="flex items-center">
              <div className="p-2 bg-primary-100 rounded-lg">
                <CalendarClockIcon className="h-6 w-6 text-primary-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-neutral-600">Deadlines</p>
                <p className="text-2xl font-bold text-neutral-900">{deadlines?.length || 0}</p>
              </div>
            </div>
          </div>
        </div>

        {/* Quick access scrollable row for deadlines*/}
        {!deadlinesLoading && deadlines && deadlines.length > 0 && (
          <div className="mb-8">
            <h2 className="text-xl font-bold text-neutral-900 mb-4">Aankomende Deadlines</h2>
            {/* Search field */}
            <div className="relative flex-1 mb-4">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-neutral-400" />
              <Input
                ref={searchInputRef}
                type="text"
                placeholder="Zoek vakken op naam, code of docent..."
                className="w-full border border-neutral-300 rounded-lg pl-10 pr-20 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                disabled={deadlinesLoading}
                value={deadlineSearchValue}
                onChange={(e) => setDeadlineSearchValue(e.target.value)}
                onKeyDown={handleSearchKeyDown}
                aria-label="Zoek deadlines"
                aria-describedby="search-shortcuts"
              />
              <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center space-x-1">
                {deadlineSearchValue && filteredDeadlines.length > 0 && (
                  <div
                    className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                    title="Druk Enter om naar eerste resultaat te gaan"
                  >
                    <CornerDownLeft className="h-3 w-3" />
                  </div>
                )}
                {!deadlineSearchValue && (
                  <div
                    className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                    title="Druk S om te zoeken"
                  >
                    <Command className="h-3 w-3 mr-0.5" />
                    <span>S</span>
                  </div>
                )}
                {deadlineSearchValue && filteredDeadlines.length > 1 && (
                  <div
                    className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                    title={`Druk 1-${Math.min(filteredDeadlines.length, 9)} om naar specifiek vak te gaan`}
                  >
                    <span>1-{Math.min(filteredDeadlines.length, 9)}</span>
                  </div>
                )}
                {deadlineSearchValue && (
                  <button
                    onClick={() => setDeadlineSearchValue('')}
                    className="p-1 text-neutral-400 hover:text-neutral-600 transition-colors rounded"
                    title="Zoekterm wissen en unfocus (Esc)"
                    aria-label="Zoekterm wissen"
                  >
                    <X className="h-3 w-3" />
                  </button>
                )}
              </div>
            </div>
            <div className="flex space-x-4 overflow-x-auto pb-2">
              {filteredDeadlines.slice(0, 5).map((deadline) => (
                <div
                  key={deadline.id}
                  className="min-w-62.5 bg-white rounded-lg border border-neutral-200 p-4 shrink-0 relative flex flex-col justify-between"
                >
                  <div>
                    <h3 className="text-lg font-medium text-neutral-900 mb-1">{deadline.title}</h3>
                    <p className="text-sm text-neutral-600 mb-2 flex items-center gap-1">
                      <span
                        className="inline-block w-3 h-3 rounded-full mr-1"
                        style={{ backgroundColor: deadline.subject.color }}
                      ></span>
                      {deadline.subject.name} ({deadline.subject.code})
                    </p>
                    <p className="text-sm text-neutral-500">
                      Inleveren op:{' '}
                      <span className="font-medium text-neutral-900">
                        {new Date(deadline.dueDate).toLocaleDateString('nl-NL', {
                          day: 'numeric',
                          month: 'long',
                          year: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </span>
                    </p>
                  </div>
                  {user?.role === 'student' ? (
                    <Link
                      to={`/deadlines/${deadline.id}/submission`}
                      className="absolute right-4 top-1/2 -translate-y-1/2 transform text-primary-600 hover:text-primary-700 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 rounded px-2 py-1 hover:bg-primary-50 transition-colors duration-300"
                      style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                    >
                      <ArrowRight className="h-5 w-5" />
                    </Link>
                  ) : (
                    <Link
                      to={`/deadlines/${deadline.id}/submissions`}
                      className="absolute right-4 top-1/2 -translate-y-1/2 transform text-primary-600 hover:text-primary-700 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 rounded px-2 py-1 hover:bg-primary-50 transition-colors duration-300"
                      style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                    >
                      <ArrowRight className="h-5 w-5" />
                    </Link>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Main Content Area */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content - Subjects */}
          <div className="lg:col-span-2">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-neutral-900">Mijn Vakken</h2>
              <div className="text-sm text-neutral-500" aria-live="polite">
                {debouncedSubjectSearchValue.trim() ? (
                  <>
                    {filteredSubjects.length} van {subjects?.length}{' '}
                    {subjects?.length === 1 ? 'vak' : 'vakken'}
                  </>
                ) : (
                  <>
                    {subjects?.length} {subjects?.length === 1 ? 'vak' : 'vakken'}
                  </>
                )}
              </div>
            </div>

            {/* Filter and search bar */}
            <div className="flex items-center gap-4 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-neutral-400" />
                <Input
                  ref={searchInputRef}
                  type="text"
                  placeholder="Zoek vakken op naam, code of docent..."
                  className="w-full border border-neutral-300 rounded-lg pl-10 pr-20 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  disabled={subjectsLoading}
                  value={subjectSearchValue}
                  onChange={(e) => setSubjectSearchValue(e.target.value)}
                  onKeyDown={handleSearchKeyDown}
                  aria-label="Zoek vakken"
                  aria-describedby="search-shortcuts"
                />
                <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center space-x-1">
                  {subjectSearchValue && filteredSubjects.length > 0 && (
                    <div
                      className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                      title="Druk Enter om naar eerste resultaat te gaan"
                    >
                      <CornerDownLeft className="h-3 w-3" />
                    </div>
                  )}
                  {!subjectSearchValue && (
                    <div
                      className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                      title="Druk S om te zoeken"
                    >
                      <Command className="h-3 w-3 mr-0.5" />
                      <span>S</span>
                    </div>
                  )}
                  {subjectSearchValue && filteredSubjects.length > 1 && (
                    <div
                      className="flex items-center px-1.5 py-0.5 bg-neutral-100 rounded text-xs text-neutral-500"
                      title={`Druk 1-${Math.min(filteredSubjects.length, 9)} om naar specifiek vak te gaan`}
                    >
                      <span>1-{Math.min(filteredSubjects.length, 9)}</span>
                    </div>
                  )}
                  {subjectSearchValue && (
                    <button
                      onClick={() => setSubjectSearchValue('')}
                      className="p-1 text-neutral-400 hover:text-neutral-600 transition-colors rounded"
                      title="Zoekterm wissen en unfocus (Esc)"
                      aria-label="Zoekterm wissen"
                    >
                      <X className="h-3 w-3" />
                    </button>
                  )}
                </div>
              </div>
              <Restricted>
                <Button
                  variant={managingSubjects ? 'outline' : 'primary'}
                  onClick={() => setManagingSubjects(!managingSubjects)}
                  className="flex items-center"
                >
                  {managingSubjects ? (
                    <>
                      <Check className="h-4 w-4 mr-2" />
                      Klaar
                    </>
                  ) : (
                    'Beheer vakken'
                  )}
                </Button>
              </Restricted>
              <div id="search-shortcuts" className="sr-only" aria-live="polite">
                {subjectSearchValue && filteredSubjects.length > 0
                  ? `Druk Enter om naar het eerste resultaat te gaan, of druk 1-${Math.min(filteredSubjects.length, 9)} om naar een specifiek vak te gaan. Escape om te wissen en unfocus.`
                  : 'Druk S om te zoeken'}
              </div>
            </div>

            {subjectsLoading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {[...Array(4)].map((_, i) => (
                  <div
                    key={i}
                    className="bg-white rounded-lg border border-neutral-200 p-6 animate-pulse"
                  >
                    <div className="h-4 bg-neutral-200 rounded w-3/4 mb-2"></div>
                    <div className="h-3 bg-neutral-200 rounded w-1/2 mb-4"></div>
                    <div className="h-3 bg-neutral-200 rounded w-full"></div>
                  </div>
                ))}
              </div>
            ) : subjects && subjects.length > 0 ? (
              filteredSubjects.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {filteredSubjects.map((subject, index) => (
                    <div
                      key={subject.id}
                      className={
                        index === 0 && debouncedSubjectSearchValue.trim()
                          ? 'ring-2 ring-primary-200 rounded-lg'
                          : ''
                      }
                    >
                      <SubjectCard
                        subject={subject}
                        searchTerm={debouncedSubjectSearchValue}
                        keyboardShortcut={index < 10 ? index : undefined}
                        managingSubject={managingSubjects}
                      />
                    </div>
                  ))}
                </div>
              ) : (
                <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
                  <Search className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
                  <h3 className="text-lg font-medium text-neutral-900 mb-2">
                    Geen vakken gevonden
                  </h3>
                  <p className="text-neutral-600 mb-4">
                    Er zijn geen vakken gevonden die overeenkomen met "{debouncedSubjectSearchValue}
                    ".
                  </p>
                  <button
                    onClick={() => setSubjectSearchValue('')}
                    className="text-primary-600 hover:text-primary-700 font-medium focus:outline-none focus:ring-2 focus:ring-primary-500 rounded px-2 py-1"
                    aria-label="Zoekterm wissen"
                  >
                    Zoekterm wissen
                  </button>
                </div>
              )
            ) : (
              <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
                <BookOpen className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-neutral-900 mb-2">Geen vakken gevonden</h3>
                <p className="text-neutral-600">
                  Je bent nog niet ingeschreven voor vakken. Neem contact op met je docent of
                  administratie.
                </p>
              </div>
            )}
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            <QuickActions />
          </div>
        </div>
      </div>
    </div>
  );
}
