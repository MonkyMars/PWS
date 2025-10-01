import { useParams, useNavigate } from 'react-router';
import { SubjectDetail } from '~/components/subjects/subject-detail';
import { ProtectedRoute } from '~/components';
import { useEffect } from 'react';

export function meta() {
  return [
    { title: 'Vak Details | PWS ELO' },
    { name: 'description', content: 'Bekijk details van je vak' },
  ];
}

export default function SubjectDetailPage() {
  const { subjectId } = useParams();
  const navigate = useNavigate();

  // Handle navigation when subjectId is missing
  useEffect(() => {
    if (!subjectId) {
      navigate('/dashboard', { replace: true });
    }
  }, [subjectId, navigate]);

  // Don't render content if no subjectId (navigation will happen via useEffect)
  if (!subjectId) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <ProtectedRoute>
      <SubjectDetail subjectId={subjectId} />
    </ProtectedRoute>
  );
}
