import { Link } from 'react-router';
import { ChevronRight, MinusCircle } from 'lucide-react';
import type { Subject } from '~/types';

interface SubjectCardProps {
  subject: Subject;
  searchTerm?: string;
  managingSubject?: boolean;
  keyboardShortcut?: number;
}

export function SubjectCard({
  subject,
  searchTerm,
  managingSubject,
  keyboardShortcut,
}: SubjectCardProps) {
  const getSubjectColor = () => {
    return subject.color;
  };

  const highlightText = (text: string, highlight?: string) => {
    if (!highlight || !highlight.trim()) return text;

    const parts = text.split(new RegExp(`(${highlight})`, 'gi'));
    return parts.map((part, index) =>
      part.toLowerCase() === highlight.toLowerCase() ? (
        <mark key={index} className="bg-yellow-200 text-yellow-900 rounded">
          {part}
        </mark>
      ) : (
        part
      )
    );
  };

  const cardContent = (
    <div className="p-6 relative">
      {/* Keyboard shortcut badge */}
      {keyboardShortcut !== undefined && !managingSubject && (
        <div className="absolute bottom-2 right-2 flex items-center justify-center w-5 h-5 bg-neutral-400/90 text-white text-xs rounded font-medium shadow-sm">
          {keyboardShortcut}
        </div>
      )}
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-3">
          <div className="w-4 h-4 rounded-full" style={{ backgroundColor: getSubjectColor() }} />
          <div>
            <h3 className="font-semibold text-neutral-900 group-hover:text-primary-600 transition-colors">
              {highlightText(subject.name, searchTerm)}
            </h3>
            <p className="text-sm text-neutral-500">{highlightText(subject.code, searchTerm)}</p>
          </div>
        </div>
        {managingSubject ? (
          <button
            className="p-1 hover:bg-neutral-100 rounded transition-colors"
            aria-label="Verwijder vak"
          >
            <MinusCircle className="text-secondary-500 w-5 h-5 cursor-pointer" />
          </button>
        ) : (
          <ChevronRight className="h-5 w-5 text-neutral-400 group-hover:text-primary-600 transition-colors" />
        )}
      </div>
    </div>
  );

  if (managingSubject) {
    return (
      <div className="bg-white rounded-lg border border-neutral-200 hover:border-neutral-300 transition-all duration-200 hover:shadow-md">
        {cardContent}
      </div>
    );
  }

  return (
    <Link
      to={`/subjects/${subject.id}`}
      className="block bg-white rounded-lg border border-neutral-200 hover:border-neutral-300 transition-all duration-200 hover:shadow-md group"
    >
      {cardContent}
    </Link>
  );
}
