import React, { useState } from 'react';
import {
  UploadButton,
  SingleFileUploadButton,
  MultipleFileUploadButton,
  ImageUploadButton,
  DocumentUploadButton,
} from './upload-button';
import type { UploadResult, FileUploadProgress, File } from '~/types';

interface UploadExampleProps {
  subjectId: string;
  onFilesUploaded?: (files: File[]) => void;
  className?: string;
}

type ButtonType = 'single' | 'multiple' | 'image' | 'document' | 'custom';

export const UploadExample: React.FC<UploadExampleProps> = ({
  subjectId,
  onFilesUploaded,
  className = '',
}) => {
  const [buttonType, setButtonType] = useState<ButtonType>('single');
  const [uploadedFiles, setUploadedFiles] = useState<File[]>([]);
  const [isUploading, setIsUploading] = useState(false);

  const handleUploadComplete = (result: UploadResult) => {
    if (result.success && result.files) {
      const newFiles = result.files;
      setUploadedFiles(prev => [...prev, ...newFiles]);
      onFilesUploaded?.(newFiles);
    }
  };

  const handleUploadStart = () => {
    setIsUploading(true);
  };

  const handleProgress = (progress: FileUploadProgress[]) => {
    // All files completed
    if (progress.every(p => p.status === 'completed' || p.status === 'error')) {
      setIsUploading(false);
    }
  };

  const renderUploadButton = () => {
    const commonProps = {
      config: { subjectId },
      onUploadComplete: handleUploadComplete,
      onUploadStart: handleUploadStart,
      onProgress: handleProgress,
      disabled: isUploading,
    };

    switch (buttonType) {
      case 'single':
        return <SingleFileUploadButton {...commonProps} />;

      case 'multiple':
        return <MultipleFileUploadButton {...commonProps} />;

      case 'image':
        return (
          <ImageUploadButton
            {...commonProps}
            buttonText="Upload Images Only"
            variant="secondary"
          />
        );

      case 'document':
        return (
          <DocumentUploadButton
            {...commonProps}
            buttonText="Upload Documents"
            variant="outline"
          />
        );

      case 'custom':
        return (
          <UploadButton
            {...commonProps}
            config={{
              subjectId,
              allowMultiple: true,
              maxFiles: 3,
              acceptedMimeTypes: ['application/pdf', 'image/*'],
            }}
            buttonText="🎯 Custom Upload (PDF + Images, max 3)"
            variant="primary"
            size="lg"
          />
        );

      default:
        return null;
    }
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Button Type Selector */}
      <div className="bg-white p-4 rounded-lg border">
        <h3 className="text-lg font-semibold mb-3">Choose Upload Button Type</h3>
        <div className="flex flex-wrap gap-2">
          {[
            { type: 'single' as const, label: 'Single File' },
            { type: 'multiple' as const, label: 'Multiple Files' },
            { type: 'image' as const, label: 'Images Only' },
            { type: 'document' as const, label: 'Documents Only' },
            { type: 'custom' as const, label: 'Custom Config' },
          ].map(({ type, label }) => (
            <button
              key={type}
              onClick={() => setButtonType(type)}
              className={`px-3 py-2 text-sm rounded-md font-medium transition-colors ${
                buttonType === type
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </div>

      {/* Current Upload Button */}
      <div className="bg-white p-6 rounded-lg border">
        <h3 className="text-lg font-semibold mb-3">
          Current Button: {buttonType.charAt(0).toUpperCase() + buttonType.slice(1)}
        </h3>
        <div className="flex justify-center">
          {renderUploadButton()}
        </div>
      </div>

      {/* Uploaded Files Display */}
      {uploadedFiles.length > 0 && (
        <div className="bg-white p-6 rounded-lg border">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">
              Uploaded Files ({uploadedFiles.length})
            </h3>
            <button
              onClick={() => setUploadedFiles([])}
              className="px-3 py-1 text-sm bg-red-100 text-red-700 rounded-md hover:bg-red-200 transition-colors"
            >
              Clear All
            </button>
          </div>
          <div className="space-y-2 max-h-60 overflow-y-auto">
            {uploadedFiles.map((file, index) => (
              <div
                key={index}
                className="flex items-center gap-3 p-3 bg-gray-50 rounded-md"
              >
                <div className="text-xl">
                  {file.mime_type.startsWith('image/') ? '🖼️' :
                   file.mime_type.includes('pdf') ? '📄' :
                   file.mime_type.includes('video') ? '🎥' :
                   file.mime_type.includes('word') ? '📝' : '📋'}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium text-gray-900 truncate">{file.name}</p>
                  <p className="text-xs text-gray-500">{file.mime_type}</p>
                </div>
                <a
                  href={file.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
                >
                  View
                </a>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Code Example */}
      <div className="bg-gray-50 p-4 rounded-lg border">
        <h3 className="text-sm font-semibold text-gray-900 mb-2">
          Current Implementation:
        </h3>
        <pre className="bg-white p-3 rounded text-xs overflow-x-auto border">
          <code>
{buttonType === 'single' && `<SingleFileUploadButton
  config={{ subjectId: "${subjectId}" }}
  onUploadComplete={handleUploadComplete}
/>`}
{buttonType === 'multiple' && `<MultipleFileUploadButton
  config={{ subjectId: "${subjectId}" }}
  onUploadComplete={handleUploadComplete}
/>`}
{buttonType === 'image' && `<ImageUploadButton
  config={{ subjectId: "${subjectId}" }}
  buttonText="Upload Images Only"
  variant="secondary"
  onUploadComplete={handleUploadComplete}
/>`}
{buttonType === 'document' && `<DocumentUploadButton
  config={{ subjectId: "${subjectId}" }}
  buttonText="Upload Documents"
  variant="outline"
  onUploadComplete={handleUploadComplete}
/>`}
{buttonType === 'custom' && `<UploadButton
  config={{
    subjectId: "${subjectId}",
    allowMultiple: true,
    maxFiles: 3,
    acceptedMimeTypes: ['application/pdf', 'image/*'],
  }}
  buttonText="🎯 Custom Upload (PDF + Images, max 3)"
  variant="primary"
  size="lg"
  onUploadComplete={handleUploadComplete}
/>`}
          </code>
        </pre>
        <p className="text-xs text-gray-600 mt-2">
          💡 Simply change the button component to swap functionality - no other changes needed!
        </p>
      </div>
    </div>
  );
};

// Higher-order component for creating domain-specific upload buttons
export const createDomainUploadButton = (
  domain: string,
  defaultConfig: Partial<Parameters<typeof UploadButton>[0]['config']>
) => {
  return (props: Omit<Parameters<typeof UploadButton>[0], 'config'> & {
    config?: Partial<Parameters<typeof UploadButton>[0]['config']>;
  }) => {
    const mergedConfig = {
      ...defaultConfig,
      ...props.config,
    } as Parameters<typeof UploadButton>[0]['config'];

    return (
      <div className="space-y-2">
        <div className="text-xs text-gray-500 font-medium uppercase tracking-wide">
          {domain} Upload
        </div>
        <UploadButton {...props} config={mergedConfig} />
      </div>
    );
  };
};

// Example domain-specific buttons
export const HomeworkUploadButton = createDomainUploadButton('Homework', {
  allowMultiple: true,
  maxFiles: 5,
  acceptedMimeTypes: ['application/pdf', 'image/*', 'text/*'],
});

export const PresentationUploadButton = createDomainUploadButton('Presentation', {
  allowMultiple: false,
  acceptedMimeTypes: [
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'application/pdf',
  ],
});

export const AssignmentUploadButton = createDomainUploadButton('Assignment', {
  allowMultiple: true,
  maxFiles: 10,
});

export default UploadExample;
