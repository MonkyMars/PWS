import { useQuery } from '@tanstack/react-query';
import { apiClient } from '~/lib/api-client';

/**
 * Represents a student's submission for a deadline.
 */
export interface StudentSubmission {
  id: string;
  deadlineId: string;
  studentId: string;
  fileIds: string[];
  message: string;
  createdAt: string;
  updatedAt: string;
  isLate: boolean;
  isUpdated: boolean;
}

/**
 * Fetches the current student's submission for a specific deadline.
 * Calls GET /deadlines/:id/submission
 *
 * @param deadlineId - The ID of the deadline to fetch the submission for.
 */
export function useStudentSubmission(deadlineId: string | undefined) {
  return useQuery({
    queryKey: ['student-submission', deadlineId],
    enabled: !!deadlineId,
    queryFn: async (): Promise<StudentSubmission | null> => {
      if (!deadlineId) return null;
      const response = await apiClient.get<StudentSubmission>(
        `/deadlines/${deadlineId}/submission`
      );
      if (!response.success) {
        throw new Error(response.message || 'Fout bij ophalen inzending');
      }
      return response.data ?? null;
    },
  });
}
