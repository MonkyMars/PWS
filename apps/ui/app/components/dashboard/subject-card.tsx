import { Link } from 'react-router';
import { ChevronRight } from 'lucide-react';
import type { Subject, Teacher } from '~/types';

interface SubjectCardProps {
  subject: Subject;
  teachers: Teacher[] | undefined;
  searchTerm?: string;
}

export function SubjectCard({ subject, searchTerm, teachers }: SubjectCardProps) {
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

  return (
    <Link
      to={`/subjects/${subject.id}`}
      className="block bg-white rounded-lg border border-neutral-200 hover:border-neutral-300 transition-all duration-200 hover:shadow-md group"
    >
      <div className="p-6">
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
          <ChevronRight className="h-5 w-5 text-neutral-400 group-hover:text-primary-600 transition-colors" />
        </div>

        {/* Teachers */}
        {teachers && teachers?.length > 0 && (
          <div className="flex justify-end mr-2">
            <div className="flex items-center">
              {teachers
                .filter((teacher) => teacher.subjectId === subject.id)
                .slice(0, 3)
                .map((teacher, index) => {
                  const colors = [
                    'bg-yellow-500',
                    'bg-purple-500',
                    'bg-green-500',
                    'bg-red-500',
                    'bg-blue-500',
                  ];
                  const randomColor = colors[Math.floor(Math.random() * colors.length)];
                  return (
                    <div
                      key={teacher.id || index}
                      className={`w-6 h-6 rounded-full ${randomColor} text-white flex items-center justify-center text-xs font-medium border-2 border-white`}
                      style={{ marginLeft: index > 0 ? '-0.5rem' : '0' }}
                    >
                      {teacher.username.charAt(0).toUpperCase()}
                    </div>
                  );
                })}
              {teachers.filter((teacher) => teacher.subjectId === subject.id).length > 3 && (
                <div className="ml-2 text-sm text-neutral-600 font-medium">
                  +{teachers.filter((teacher) => teacher.subjectId === subject.id).length - 3}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </Link>
  );
}
