/**
 * Drive file information from Google Drive API
 */
export interface DriveFile {
  file_id: string;
  name: string;
  mime_type: string;
}

/**
 * File entity as stored in the database
 */
export interface File {
  id: string;
  file_id: string;
  name: string;
  mime_type: string;
  subject_id: string;
  uploaded_by: string;
  url: string;
  created_at: string;
  updated_at: string;
}

/**
 * Upload request for single file
 */
export interface UploadSingleFileRequest {
  file: DriveFile;
  subject_id: string;
}

/**
 * Upload request for multiple files
 */
export interface UploadMultipleFilesRequest {
  files: DriveFile[];
  subject_id: string;
}

/**
 * File upload progress information
 */
export interface FileUploadProgress {
  fileName: string;
  progress: number;
  status: 'pending' | 'uploading' | 'completed' | 'error';
  error?: string;
}

/**
 * Upload configuration options
 */
export interface UploadConfig {
  subjectId: string;
  allowMultiple?: boolean;
  maxFiles?: number;
  acceptedMimeTypes?: string[];
  maxFileSize?: number; // in bytes
}

/**
 * File upload result
 */
export interface UploadResult {
  success: boolean;
  files?: File[];
  errors?: string[];
}
