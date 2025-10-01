/**
 * Standard API response wrapper
 */
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  errors?: ApiError[];
}

/**
 * API error interface
 */
export interface ApiError {
  field?: string;
  message: string;
  code?: string;
}

/**
 * Pagination parameters for list requests
 */
export interface PaginationParams {
  page?: number;
  limit?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

/**
 * Paginated response wrapper
 */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

/**
 * File upload progress information
 */
export interface UploadProgress {
  loaded: number;
  total: number;
  percentage: number;
}

/**
 * Filter options for subjects
 */
export interface SubjectFilters {
  teacherId?: string;
  isActive?: boolean;
  search?: string;
}

/**
 * Filter options for announcements
 */
export interface AnnouncementFilters {
  subjectId?: string;
  priority?: string;
  isSticky?: boolean;
  authorId?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * Filter options for files
 */
export interface FileFilters {
  subjectId?: string;
  folderId?: string;
  category?: string;
  mimeType?: string;
  uploaderId?: string;
  isPublic?: boolean;
}
