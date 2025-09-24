import { apiClient } from './api-client';
import type {
  DriveFile,
  File,
  UploadSingleFileRequest,
  UploadMultipleFilesRequest,
  UploadResult,
  FileUploadProgress,
  ApiResponse,
} from '~/types';

/**
 * File upload service for handling Google Drive file picker and API uploads
 */
export class FileUploadService {
  /**
   * Upload a single file to the API
   */
  static async uploadSingle(file: DriveFile, subjectId: string): Promise<ApiResponse<File>> {
    const request = {
      file,
      subject_id: subjectId,
    };

    return apiClient.post<File>('/files/upload/single', request);
  }

  /**
   * Upload multiple files to the API
   */
  static async uploadMultiple(files: DriveFile[], subjectId: string): Promise<ApiResponse<File[]>> {
    const request: UploadMultipleFilesRequest = {
      files,
      subject_id: subjectId,
    };

    return apiClient.post<File[]>('/files/upload/multiple', request);
  }

  /**
   * Get Google Drive access token for current user
   */
  static async getGoogleAccessToken(): Promise<string> {
    const response = await apiClient.get<{ access_token: string }>('/auth/google/access-token');

    if (!response.success || !response.data) {
      throw new Error(response.message || 'Failed to get Google access token');
    }

    return response.data.access_token;
  }

  /**
   * Check if user has linked their Google account
   */
  static async checkGoogleLinkStatus(): Promise<boolean> {
    const response = await apiClient.get<{ linked: boolean }>('/auth/google/status');
    return response.success && response.data?.linked === true;
  }

  /**
   * Get Google OAuth URL for linking account
   */
  static async getGoogleAuthURL(): Promise<string> {
    const response = await apiClient.get<{ auth_url: string }>('/auth/google/url');

    if (!response.success || !response.data) {
      throw new Error(response.message || 'Failed to get Google auth URL');
    }

    return response.data.auth_url;
  }

  /**
   * Initialize Google Drive Picker
   */
  static initializePicker(
    accessToken: string,
    onSelection: (files: DriveFile[]) => void,
    allowMultiple: boolean = false
  ): void {
    const pickerCallback = (data: any) => {
      if (data.action === google.picker.Action.PICKED) {
        const selectedFiles: DriveFile[] = data.docs.map((doc: any) => ({
          file_id: doc.id,
          name: doc.name,
          mime_type: doc.mimeType,
        }));

        onSelection(selectedFiles);
      }
    };

    let builder = new google.picker.PickerBuilder()
      .addView(google.picker.ViewId.DOCS)
      .addView(google.picker.ViewId.DOCS_IMAGES)
      .addView(google.picker.ViewId.DOCS_VIDEOS)
      .setOAuthToken(accessToken)
      .setDeveloperKey(import.meta.env.VITE_GOOGLE_API_KEY || '')
      .setCallback(pickerCallback)
      .setSize(1051, 650);

    if (!allowMultiple) {
      builder = builder.setSelectableMimeTypes('application/pdf,image/*,video/*,text/*');
    }

    const picker = builder.build();
    picker.setVisible(true);
  }

  /**
   * Load Google Picker API script
   */
  static loadGooglePickerScript(): Promise<void> {
    return new Promise((resolve, reject) => {
      // Check if script is already loaded
      if (window.google?.picker) {
        resolve();
        return;
      }

      // Create script element
      const script = document.createElement('script');
      script.src = 'https://apis.google.com/js/api.js';
      script.onload = () => {
        // Load picker API
        gapi.load('picker', {
          callback: () => resolve(),
          onerror: () => reject(new Error('Failed to load Google Picker API')),
        });
      };
      script.onerror = () => reject(new Error('Failed to load Google APIs script'));

      document.head.appendChild(script);
    });
  }

  /**
   * Open Google Drive Picker and handle file selection
   */
  static async openPicker(
    subjectId: string,
    allowMultiple: boolean = false,
    onProgress?: (progress: FileUploadProgress[]) => void
  ): Promise<UploadResult> {
    try {
      // Check if user has linked Google account
      const isLinked = await this.checkGoogleLinkStatus();
      if (!isLinked) {
        const authUrl = await this.getGoogleAuthURL();
        window.open(authUrl, '_blank', 'width=500,height=600');
        throw new Error('Please link your Google account first');
      }

      // Load Google Picker API
      await this.loadGooglePickerScript();

      // Get access token
      const accessToken = await this.getGoogleAccessToken();

      // Open picker and wait for selection
      return new Promise((resolve, reject) => {
        const onSelection = async (selectedFiles: DriveFile[]) => {
          try {
            const progressList: FileUploadProgress[] = selectedFiles.map((file) => ({
              fileName: file.name,
              progress: 0,
              status: 'pending' as const,
            }));

            const updateProgress = () => {
              onProgress?.(progressList);
            };

            // Update progress to uploading
            progressList.forEach((p) => (p.status = 'uploading'));
            updateProgress();

            let apiResponse: ApiResponse<File | File[]>;

            if (selectedFiles.length === 1) {
              apiResponse = await this.uploadSingle(selectedFiles[0], subjectId);
            } else {
              apiResponse = await this.uploadMultiple(selectedFiles, subjectId);
            }

            // Update progress to completed
            progressList.forEach((p) => {
              p.progress = 100;
              p.status = 'completed';
            });
            updateProgress();

            if (!apiResponse.success) {
              throw new Error(apiResponse.message || 'Upload failed');
            }

            resolve({
              success: true,
              files: Array.isArray(apiResponse.data) ? apiResponse.data : [apiResponse.data!],
            });
          } catch (error) {
            reject({
              success: false,
              errors: [error instanceof Error ? error.message : 'Unknown error occurred'],
            });
          }
        };

        this.initializePicker(accessToken, onSelection, allowMultiple);
      });
    } catch (error) {
      return {
        success: false,
        errors: [error instanceof Error ? error.message : 'Unknown error occurred'],
      };
    }
  }

  /**
   * Format file size for display
   */
  static formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  /**
   * Get file icon based on mime type
   */
  static getFileIcon(mimeType: string): string {
    if (mimeType.startsWith('image/')) return '🖼️';
    if (mimeType.startsWith('video/')) return '🎥';
    if (mimeType.startsWith('audio/')) return '🎵';
    if (mimeType.includes('pdf')) return '📄';
    if (mimeType.includes('word') || mimeType.includes('document')) return '📝';
    if (mimeType.includes('sheet') || mimeType.includes('excel')) return '📊';
    if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) return '📈';
    if (mimeType.includes('zip') || mimeType.includes('compressed')) return '📦';
    return '📋';
  }
}

// Global type declarations for Google APIs
declare global {
  interface Window {
    google?: {
      picker: {
        PickerBuilder: new () => GooglePickerBuilder;
        ViewId: {
          DOCS: string;
          DOCS_IMAGES: string;
          DOCS_VIDEOS: string;
        };
        Action: {
          PICKED: string;
        };
      };
    };
    gapi?: {
      load: (api: string, options: { callback: () => void; onerror: () => void }) => void;
    };
  }

  interface GooglePickerBuilder {
    addView(viewId: string): GooglePickerBuilder;
    setOAuthToken(token: string): GooglePickerBuilder;
    setDeveloperKey(key: string): GooglePickerBuilder;
    setCallback(callback: (data: any) => void): GooglePickerBuilder;
    setSize(width: number, height: number): GooglePickerBuilder;
    setSelectableMimeTypes(types: string): GooglePickerBuilder;
    build(): GooglePicker;
  }

  interface GooglePicker {
    setVisible(visible: boolean): void;
  }

  const google: {
    picker: {
      PickerBuilder: new () => GooglePickerBuilder;
      ViewId: {
        DOCS: string;
        DOCS_IMAGES: string;
        DOCS_VIDEOS: string;
      };
      Action: {
        PICKED: string;
      };
    };
  };

  const gapi: {
    load: (api: string, options: { callback: () => void; onerror: () => void }) => void;
  };
}
