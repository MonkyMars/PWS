import { useParams } from 'react-router';
import { SubjectDetail } from '~/components/subjects/subject-detail';
import { ProtectedRoute } from '~/components';
import { Navigate } from 'react-router';

export function meta() {
  return [
    { title: 'Vak Details | PWS ELO' },
    { name: 'description', content: 'Bekijk details van je vak' },
  ];
}

export default function SubjectDetailPage() {
  const { subjectId } = useParams();

  if (!subjectId) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <ProtectedRoute>
      <SubjectDetail subjectId={subjectId} />
    </ProtectedRoute>
  );
}
