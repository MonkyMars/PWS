import { BookOpen, Bell, FileText, Calendar, TrendingUp } from "lucide-react";
import { SubjectCard } from "./subject-card";
import { RecentActivity } from "./recent-activity";
import { QuickActions } from "./quick-actions";
import { useCurrentUser, useSubjects } from "~/hooks";

export function Dashboard() {
  const { data: user } = useCurrentUser();
  const { data: subjects, isLoading: subjectsLoading } = useSubjects();

  // if (subjectsLoading) {
  //   return (
  //     <div className="min-h-screen flex items-center justify-center">
  //       <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
  //     </div>
  //   );
  // }

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Goedemorgen";
    if (hour < 17) return "Goedemiddag";
    return "Goedenavond";
  };

  return (
    <div className="min-h-screen bg-neutral-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Welcome Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-neutral-900 mb-2">
            {getGreeting()}, {user?.username}!
          </h1>
          <p className="text-neutral-600">
            Welkom terug in je PWS ELO dashboard. Hier vind je al je vakken en
            recente activiteiten.
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
                <p className="text-2xl font-bold text-neutral-900">
                  {subjects?.length || 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 border border-neutral-200">
            <div className="flex items-center">
              <div className="p-2 bg-warning-100 rounded-lg">
                <Bell className="h-6 w-6 text-warning-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-neutral-600">
                  Nieuwe Mededelingen
                </p>
                <p className="text-2xl font-bold text-neutral-900">3</p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 border border-neutral-200">
            <div className="flex items-center">
              <div className="p-2 bg-success-100 rounded-lg">
                <FileText className="h-6 w-6 text-success-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-neutral-600">
                  Nieuwe Bestanden
                </p>
                <p className="text-2xl font-bold text-neutral-900">7</p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 border border-neutral-200">
            <div className="flex items-center">
              <div className="p-2 bg-secondary-100 rounded-lg">
                <TrendingUp className="h-6 w-6 text-secondary-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-neutral-600">
                  Voortgang
                </p>
                <p className="text-2xl font-bold text-neutral-900">89%</p>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content - Subjects */}
          <div className="lg:col-span-2">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-neutral-900">
                Mijn Vakken
              </h2>
              <div className="text-sm text-neutral-500">
                {subjects?.length} {subjects?.length === 1 ? "vak" : "vakken"}
              </div>
            </div>

            {subjects && subjects.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {subjects.map((subject) => (
                  <SubjectCard key={subject.id} subject={subject} />
                ))}
              </div>
            ) : (
              <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
                <BookOpen className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-neutral-900 mb-2">
                  Geen vakken gevonden
                </h3>
                <p className="text-neutral-600">
                  Je bent nog niet ingeschreven voor vakken. Neem contact op met
                  je docent of administratie.
                </p>
              </div>
            )}
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            <QuickActions />
            <RecentActivity />
          </div>
        </div>
      </div>
    </div>
  );
}
