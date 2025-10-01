import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '~/lib/api-client';
import type {
  Subject,
  Announcement,
  SubjectFile,
  PaginatedResponse,
  SubjectFilters,
  AnnouncementFilters,
  FileFilters,
  SubjectFolder,
} from '~/types';

/**
 * Hook to get user's subjects
 */
export function useSubjects(filters?: SubjectFilters) {
  return useQuery({
    queryKey: ['subjects', filters],
    queryFn: async (): Promise<Subject[]> => {
      const response = await apiClient.get<Subject[]>('/subjects/me', filters);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen vakken');
      }

      const data = response.data.map((file: any) => {
        return {
          createdAt: new Date(file.created_at).toISOString(),
          mimeType: file.mime_type,
          teacherId: file.teacher_id,
          teacherName: file.teacher_name,
          updatedAt: new Date(file.updated_at).toISOString(),
          ...file,
        };
      });

      const sortedData = data.sort((a, b) => a.name.localeCompare(b.name));

      return sortedData;
    },
  });
}

/**
 * Hook to get a specific subject by ID
 */
export function useSubject(subjectId: string) {
  return useQuery({
    queryKey: ['subjects', subjectId],
    queryFn: async (): Promise<Subject> => {
      const response = await apiClient.get<{
        id: string;
        name: string;
        code: string;
        color: string;
        created_at: string;
        updated_at: string;
        teacher_id: string;
        teacher_name: string;
        is_active: boolean;
      }>(`/subjects/${subjectId}`);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen vak');
      }

      const data: Subject = {
        ...response.data,
        createdAt: new Date(response.data.created_at).toISOString(),
        updatedAt: new Date(response.data.updated_at).toISOString(),
        teacherId: response.data.teacher_id,
        teacherName: response.data.teacher_name,
        isActive: response.data.is_active,
      };

      return data;
    },
    enabled: !!subjectId,
  });
}

/**
 * Hook to get announcements for a subject
 */
export function useAnnouncements(filters?: AnnouncementFilters) {
  return useQuery({
    queryKey: ['announcements', filters],
    queryFn: async (): Promise<PaginatedResponse<Announcement>> => {
      const response = await apiClient.get<PaginatedResponse<Announcement>>(
        '/announcements',
        filters
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen mededelingen');
      }

      return response.data;
    },
  });
}

/**
 * Hook to get files for a subject
 */
export function useSubjectFiles(filters?: FileFilters) {
  return useQuery({
    queryKey: ['files', filters],
    queryFn: async (): Promise<PaginatedResponse<SubjectFile>> => {
      if (!filters?.subjectId) {
        throw new Error('subjectId is required to fetch files');
      }
      if (!filters.folderId) {
        filters.folderId = filters.subjectId;
      }
      const path = `/files/subject/${filters?.subjectId}/folder/${filters?.folderId}`;
      const response = await apiClient.get<PaginatedResponse<SubjectFile>>(path);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen bestanden');
      }

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen bestanden');
      }

      const data = response.data.items.map((file: any) => {
        return {
          createdAt: new Date(file.created_at).toISOString(),
          mimeType: file.mime_type,
          subjectId: file.subject_id,
          subjectName: file.subject_name,
          updatedAt: new Date(file.updated_at).toISOString(),
          folderId: file.folder_id,
          ...file,
        };
      });

      const sortedData: SubjectFile[] = data.sort((a, b) => b.createdAt.localeCompare(a.createdAt));

      return { ...response.data, items: sortedData };
    },
  });
}
/**
 * Hook to get folders for a subject
 */
export function useSubjectFolders(filters?: FileFilters) {
  return useQuery({
    queryKey: ['folders', filters],
    queryFn: async (): Promise<PaginatedResponse<SubjectFolder>> => {
      if (!filters?.subjectId) {
        throw new Error('subjectId is required to fetch files');
      }
      if (!filters.folderId) {
        filters.folderId = filters.subjectId;
      }
      const path = `/folders/subject/${filters?.subjectId}/folder/${filters?.folderId}`;
      const response = await apiClient.get<PaginatedResponse<SubjectFolder>>(path);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen bestanden');
      }

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij ophalen bestanden');
      }

      const data = response.data.items.map((folder: any) => {
        return {
          createdAt: new Date(folder.created_at).toISOString(),
          subjectId: folder.subject_id,
          updatedAt: new Date(folder.updated_at).toISOString(),
          parentId: folder.parent_id,
          uploaderId: folder.uploader_id,
          ...folder,
        };
      });

      const sortedData: SubjectFolder[] = data.sort((a, b) =>
        b.createdAt.localeCompare(a.createdAt)
      );

      return { ...response.data, items: sortedData };
    },
  });
}

/**
 * Hook to upload a file
 */
export function useUploadFile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      file,
      subjectId,
      description,
      category,
      onProgress,
    }: {
      file: File;
      subjectId: string;
      description?: string;
      category?: string;
      onProgress?: (progress: number) => void;
    }): Promise<SubjectFile> => {
      const response = await apiClient.uploadFile<SubjectFile>(
        '/files/upload',
        file,
        { subjectId, description, category },
        onProgress
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Fout bij uploaden bestand');
      }

      return response.data;
    },
    onSuccess: (data) => {
      // Invalidate file queries to refresh the list
      queryClient.invalidateQueries({ queryKey: ['files'] });

      // Update subject data if applicable
      queryClient.invalidateQueries({ queryKey: ['subjects', data.subjectId] });
    },
  });
}

/**
 * Hook to delete a file
 */
export function useDeleteFile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (fileId: string): Promise<void> => {
      const response = await apiClient.delete(`/files/${fileId}`);

      if (!response.success) {
        throw new Error(response.message || 'Fout bij verwijderen bestand');
      }
    },
    onSuccess: () => {
      // Invalidate file queries to refresh the list
      queryClient.invalidateQueries({ queryKey: ['files'] });
      queryClient.invalidateQueries({ queryKey: ['subjects'] });
    },
  });
}
