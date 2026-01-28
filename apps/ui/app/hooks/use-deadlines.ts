import { useQuery } from '@tanstack/react-query';
import type { UUID } from 'crypto';
import { apiClient } from '~/lib/api-client';
import type { ApiResponse } from '~/types';
import type { Deadline, DeadlineFilters } from '~/types/deadlines';

/**
 * Hook to get user's deadlines
 */
export function useDeadlines(filters?: DeadlineFilters) {
	return useQuery({
		queryKey: ['deadlines', filters],
		queryFn: async (): Promise<Deadline[]> => {
			const response = await apiClient.get<Deadline[]>('/deadlines/me', {
				date_from: filters?.dateFrom,
				date_to: filters?.dateTo,
				subject_id: filters?.subjectId,
				owner_id: filters?.ownerId,
			});

			if (!response.success || !response.data) {
				throw new Error(response.message || 'Fout bij ophalen inlever opdrachten');
			}

			const data: Deadline[] = response.data.map((deadline: any) => {
				return {
					createdAt: new Date(deadline.created_at).toISOString(),
					dueDate: new Date(deadline.due_date).toISOString(),
					ownerId: deadline.owner_id,
					subjectId: deadline.subject_id,
					...deadline,
				};
			});

			const sortedData = data.sort((a, b) => a.dueDate.localeCompare(b.dueDate));

			return sortedData;
		},
	});
}

export function useDeleteDeadlines(deadlineId: UUID) {
	return useQuery({
		queryKey: ['delete-deadline', deadlineId],
		queryFn: async (): Promise<{ success: boolean; message: string }> => {
			const response = await apiClient.delete<{ success: boolean; message: string }>(`/deadlines/${deadlineId}`);

			if (!response.success) {
				throw new Error(response.message || 'Fout bij verwijderen inlever opdracht');
			}

			return {
				success: response.success,
				message: response.message ?? 'Verwijderen gelukt',
			};
		},
	});
}

export function useUpdateDeadline(deadlineId: UUID, updatedData: Partial<Deadline>) {
	return useQuery({
		queryKey: ['update-deadline', deadlineId, updatedData],
		queryFn: async (): Promise<Deadline> => {
			const response = await apiClient.put<Deadline>(`/deadlines/${deadlineId}`, {
				...updatedData,
				due_date: updatedData.dueDate,
				created_at: updatedData.createdAt,
				owner_id: updatedData.ownerId,
				subject_id: updatedData.subjectId,
			});

			if (!response.success || !response.data) {
				throw new Error(response.message || 'Fout bij bijwerken inlever opdracht');
			}

			const deadline: Deadline = response.data;

			return deadline;
		},
	});
}