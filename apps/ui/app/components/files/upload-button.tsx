import React, { useState, useCallback } from 'react';
import { Upload, Link, CheckCircle, AlertCircle, Loader2, X } from 'lucide-react';
import { Button } from '../ui/button';
import { FileUploadService } from '~/lib/file-upload';
import type { FileUploadProgress, UploadResult, UploadConfig } from '~/types';

interface UploadButtonProps {
  /** Configuration for the upload */
  config: UploadConfig;
  /** Callback when upload is completed */
  onUploadComplete?: (result: UploadResult) => void;
  /** Callback when upload starts */
  onUploadStart?: () => void;
  /** Callback for progress updates */
  onProgress?: (progress: FileUploadProgress[]) => void;
  /** Custom button text */
  buttonText?: string;
  /** Button variant */
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  /** Button size */
  size?: 'sm' | 'md' | 'lg';
  /** Whether the button is disabled */
  disabled?: boolean;
  /** Additional CSS classes */
  className?: string;
}

export const UploadButton: React.FC<UploadButtonProps> = ({
  config,
  onUploadComplete,
  onUploadStart,
  onProgress,
  buttonText,
  variant = 'primary',
  size = 'md',
  disabled = false,
  className = '',
}) => {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState<FileUploadProgress[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isGoogleLinked, setIsGoogleLinked] = useState<boolean | null>(null);
  const [isLinking, setIsLinking] = useState(false);

  // Check Google link status on component mount
  React.useEffect(() => {
    const checkStatus = async () => {
      try {
        const linked = await FileUploadService.checkGoogleLinkStatus();
        setIsGoogleLinked(linked);
      } catch (err) {
        console.error('Failed to check Google link status:', err);
        setIsGoogleLinked(false);
      }
    };

    checkStatus();
  }, []);

  const handleUpload = useCallback(async () => {
    setError(null);
    setIsUploading(true);
    setIsLinking(false);
    setUploadProgress([]);

    try {
      onUploadStart?.();

      const result = await FileUploadService.openPicker(
        // Not sensitive since this is just for demo purposes
        '19601f2a-796f-4e01-a2d7-c9949daa6505', // Hard coded uuid because subject logic is not implemented yet
        {
          allowMultiple: config.allowMultiple,
          acceptedMimeTypes: config.acceptedMimeTypes,
        },
        (progress) => {
          setUploadProgress(progress);
          onProgress?.(progress);
        },
        () => {
          // OAuth flow is starting
          setIsLinking(true);
        }
      );

      if (result.success) {
        // Update Google linked status after successful upload
        setIsGoogleLinked(true);
        onUploadComplete?.(result);
      } else {
        setError(result.errors?.[0] || 'Upload failed');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Upload failed';
      setError(errorMessage);

      // If error mentions linking, it might be an OAuth issue
      if (
        errorMessage.toLowerCase().includes('link') ||
        errorMessage.toLowerCase().includes('oauth')
      ) {
        setIsGoogleLinked(false);
      }
    } finally {
      setIsUploading(false);
      setIsLinking(false);
    }
  }, [config, onUploadComplete, onUploadStart, onProgress]);

  const handleLinkGoogle = useCallback(async () => {
    try {
      const authUrl = await FileUploadService.getGoogleAuthURL();
      const popup = window.open(
        authUrl,
        'google-oauth',
        'width=500,height=600,scrollbars=yes,resizable=yes'
      );

      // Listen for popup close to check status
      const checkClosed = setInterval(() => {
        if (popup?.closed) {
          clearInterval(checkClosed);
          // Recheck link status after popup closes
          setTimeout(async () => {
            const linked = await FileUploadService.checkGoogleLinkStatus();
            setIsGoogleLinked(linked);
          }, 1000);
        }
      }, 1000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to open Google OAuth');
    }
  }, []);

  // Show loading if checking status or linking in progress
  if (isGoogleLinked === null || isLinking) {
    const loadingText =
      isGoogleLinked === null
        ? 'Checking...'
        : isLinking
          ? 'Linking Google Drive...'
          : 'Loading...';

    return (
      <Button variant={variant} size={size} disabled className={className}>
        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
        {loadingText}
      </Button>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      <Button
        variant={variant}
        size={size}
        onClick={handleUpload}
        disabled={disabled || isUploading}
        className={className}
      >
        {isUploading ? (
          <>
            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            {isLinking ? 'Linking Google Drive...' : 'Uploading...'}
          </>
        ) : (
          <>
            <Upload className="w-4 h-4 mr-2" />
            {buttonText || (config.allowMultiple ? 'Upload Files' : 'Upload File')}
          </>
        )}
      </Button>

      {/* Progress indicator */}
      {uploadProgress.length > 0 && (
        <div className="space-y-2">
          {uploadProgress.map((progress, index) => (
            <div key={index} className="flex items-center gap-2 text-sm">
              {progress.status === 'completed' ? (
                <CheckCircle className="w-4 h-4 text-green-500" />
              ) : progress.status === 'error' ? (
                <AlertCircle className="w-4 h-4 text-red-500" />
              ) : (
                <Loader2 className="w-4 h-4 animate-spin text-blue-500" />
              )}
              <span className="flex-1 truncate">{progress.fileName}</span>
              {progress.status === 'uploading' && (
                <span className="text-gray-500">{Math.round(progress.progress)}%</span>
              )}
              {progress.error && <span className="text-red-500 text-xs">{progress.error}</span>}
            </div>
          ))}
        </div>
      )}

      {/* Error message with manual link option */}
      {error && (
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-2 text-sm text-red-600">
            <AlertCircle className="w-4 h-4" />
            {error}
            <button
              onClick={() => setError(null)}
              className="ml-auto text-gray-400 hover:text-gray-600"
            >
              <X className="w-3 h-3" />
            </button>
          </div>
          {/* Show manual link option if error is related to OAuth/linking and not already linked */}
          {(error.toLowerCase().includes('link') ||
            error.toLowerCase().includes('oauth') ||
            error.toLowerCase().includes('popup')) &&
            isGoogleLinked === false && (
              <Button variant="outline" size="sm" onClick={handleLinkGoogle} className="text-xs">
                <Link className="w-3 h-3 mr-1" />
                Try Manual Link
              </Button>
            )}
        </div>
      )}
    </div>
  );
};

// Higher-order component for easy swapping
export const createUploadButton = (defaultConfig: Partial<UploadConfig>) => {
  return (props: Omit<UploadButtonProps, 'config'> & { config?: Partial<UploadConfig> }) => {
    const mergedConfig = { ...defaultConfig, ...props.config } as UploadConfig;
    return <UploadButton {...props} config={mergedConfig} />;
  };
};

// Pre-configured upload buttons for common use cases
export const SingleFileUploadButton = createUploadButton({
  allowMultiple: false,
  maxFiles: 1,
});

export const MultipleFileUploadButton = createUploadButton({
  allowMultiple: true,
  maxFiles: 10,
});

export const ImageUploadButton = createUploadButton({
  allowMultiple: false,
  acceptedMimeTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
});

export const DocumentUploadButton = createUploadButton({
  allowMultiple: true,
  acceptedMimeTypes: [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'text/plain',
  ],
});

export default UploadButton;
