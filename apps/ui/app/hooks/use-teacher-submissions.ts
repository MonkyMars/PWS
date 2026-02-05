import { useQuery } from '@tanstack/react-query';
import { apiClient } from '~/lib/api-client';

export interface TeacherSubmission {
  id: string;
  deadlineId: string;
  studentId: string;
  fileIds: string[];
  message: string;
  createdAt: string;
  updatedAt: string;
  isLate: boolean;
  isUpdated: boolean;
  student: {
    id: string;
    name: string;
    email: string;
  };
  deadline: {
    id: string;
    title: string;
    dueDate: string;
    subject: {
      id: string;
      name: string;
      code: string;
    };
  };
}

/**
 * Fetch all student submissions for deadlines the teacher/admin can view.
 * This calls the backend endpoint: GET /deadlines/:id/submissions
 *
 * Optionally, you can pass a deadlineId to filter for a specific deadline.
 */
export function useTeacherSubmissions(deadlineId?: string) {
  return useQuery({
    queryKey: ['teacher-submissions', deadlineId],
    queryFn: async (): Promise<TeacherSubmission[]> => {
      let endpoint = '/deadlines/submissions';
      if (deadlineId) {
        endpoint = `/deadlines/${deadlineId}/submissions`;
      }
      const response = await apiClient.get<TeacherSubmission[]>(endpoint);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen ingeleverde opdrachten');
      }

      // Optionally, sort by createdAt descending
      return response.data.sort((a, b) => b.createdAt.localeCompare(a.createdAt));
    },
  });
}
