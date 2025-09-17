import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "~/lib/api-client";
import type {
  Subject,
  SubjectWithDetails,
  Announcement,
  SubjectFile,
  PaginatedResponse,
  SubjectFilters,
  AnnouncementFilters,
  FileFilters,
} from "~/types";

/**
 * Hook to get user's subjects
 */
export function useSubjects(filters?: SubjectFilters) {
  return useQuery({
    queryKey: ["subjects", filters],
    queryFn: async (): Promise<SubjectWithDetails[]> => {
      const response = await apiClient.get<SubjectWithDetails[]>(
        "/subjects",
        filters
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || "Fout bij ophalen vakken");
      }

      return response.data;
    },
  });
}

/**
 * Hook to get a specific subject by ID
 */
export function useSubject(subjectId: string) {
  return useQuery({
    queryKey: ["subjects", subjectId],
    queryFn: async (): Promise<Subject> => {
      const response = await apiClient.get<Subject>(`/subjects/${subjectId}`);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Fout bij ophalen vak");
      }

      return response.data;
    },
    enabled: !!subjectId,
  });
}

/**
 * Hook to get announcements for a subject
 */
export function useAnnouncements(filters?: AnnouncementFilters) {
  return useQuery({
    queryKey: ["announcements", filters],
    queryFn: async (): Promise<PaginatedResponse<Announcement>> => {
      const response = await apiClient.get<PaginatedResponse<Announcement>>(
        "/announcements",
        filters
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || "Fout bij ophalen mededelingen");
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
    queryKey: ["files", filters],
    queryFn: async (): Promise<PaginatedResponse<SubjectFile>> => {
      const response = await apiClient.get<PaginatedResponse<SubjectFile>>(
        "/files",
        filters
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || "Fout bij ophalen bestanden");
      }

      return response.data;
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
        "/files/upload",
        file,
        { subjectId, description, category },
        onProgress
      );

      if (!response.success || !response.data) {
        throw new Error(response.message || "Fout bij uploaden bestand");
      }

      return response.data;
    },
    onSuccess: (data) => {
      // Invalidate file queries to refresh the list
      queryClient.invalidateQueries({ queryKey: ["files"] });

      // Update subject data if applicable
      queryClient.invalidateQueries({ queryKey: ["subjects", data.subjectId] });
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
        throw new Error(response.message || "Fout bij verwijderen bestand");
      }
    },
    onSuccess: () => {
      // Invalidate file queries to refresh the list
      queryClient.invalidateQueries({ queryKey: ["files"] });
      queryClient.invalidateQueries({ queryKey: ["subjects"] });
    },
  });
}
