import { Link } from 'react-router';
import { Bell, FileText, Clock, ChevronRight } from 'lucide-react';
import type { Subject } from '~/types';

interface SubjectCardProps {
  subject: Subject;
}

export function SubjectCard({ subject }: SubjectCardProps) {
  const getSubjectColor = (code: string) => {
    // Map subject codes to colors
    const colorMap: Record<string, string> = {
      WISK: 'var(--color-subject-math)',
      NATK: 'var(--color-subject-science)',
      NEDD: 'var(--color-subject-language)',
      GESCH: 'var(--color-subject-history)',
      KUNST: 'var(--color-subject-arts)',
      SPORT: 'var(--color-subject-sports)',
    };

    return colorMap[code] || 'var(--color-subject-default)';
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('nl-NL', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <Link
      to={`/subjects/${subject.id}`}
      className="block bg-white rounded-lg border border-neutral-200 hover:border-neutral-300 transition-all duration-200 hover:shadow-md group"
    >
      <div className="p-6">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div
              className="w-4 h-4 rounded-full"
              style={{ backgroundColor: getSubjectColor(subject.code) }}
            />
            <div>
              <h3 className="font-semibold text-neutral-900 group-hover:text-primary-600 transition-colors">
                {subject.name}
              </h3>
              <p className="text-sm text-neutral-500">{subject.code}</p>
            </div>
          </div>
          <ChevronRight className="h-5 w-5 text-neutral-400 group-hover:text-primary-600 transition-colors" />
        </div>

        {/* Teacher */}
        <div className="mb-4">
          <p className="text-sm text-neutral-600">
            Docent: <span className="font-medium">{subject.teacherName}</span>
          </p>
        </div>
      </div>
    </Link>
  );
}
