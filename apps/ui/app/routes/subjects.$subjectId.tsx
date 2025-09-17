import { useParams } from 'react-router';
import { SubjectDetail } from '~/components/subjects/subject-detail';
import { useCurrentUser } from '~/hooks';
import { Navigate } from 'react-router';

export function meta() {
	return [
		{ title: "Vak Details | PWS ELO" },
		{ name: "description", content: "Bekijk details van je vak" },
	];
}

export default function SubjectDetailPage() {
	const { subjectId } = useParams();
	const { data: user, isLoading } = useCurrentUser();

	if (isLoading) {
		return (
			<div className="min-h-screen flex items-center justify-center">
				<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
			</div>
		);
	}

	if (!user) {
		return <Navigate to="/login" replace />;
	}

	if (!subjectId) {
		return <Navigate to="/dashboard" replace />;
	}

	return <SubjectDetail subjectId={subjectId} />;
}