import React, { useState } from 'react';
import type { Route } from './+types/test-upload';
import {
  UploadButton,
  SingleFileUploadButton,
  MultipleFileUploadButton,
  ImageUploadButton,
  DocumentUploadButton,
} from '~/components/files/upload-button';
import type { UploadResult, FileUploadProgress, File } from '~/types';

export function meta() {
  return [
    { title: 'Test Upload - PWS' },
    { name: 'description', content: 'Test file upload functionality' },
  ];
}

export default function TestUpload() {
  const [uploadResults, setUploadResults] = useState<File[]>([]);
  const [currentProgress, setCurrentProgress] = useState<FileUploadProgress[]>([]);
  const [selectedSubjectId] = useState('19601f2a-796f-4e01-a2d7-c9949daa6505'); // Mock subject ID

  const handleUploadComplete = (result: UploadResult) => {
    if (result.success && result.files) {
      setUploadResults((prev) => [...prev, ...result.files!]);
      console.log('Upload completed:', result.files);
    } else {
      console.error('Upload failed:', result.errors);
    }
  };

  const handleUploadStart = () => {
    console.log('Upload started');
    setCurrentProgress([]);
  };

  const handleProgress = (progress: FileUploadProgress[]) => {
    setCurrentProgress(progress);
    console.log('Progress:', progress);
  };

  const clearResults = () => {
    setUploadResults([]);
    setCurrentProgress([]);
  };

  return (
    <div className="container mx-auto bg-gray-100 px-4 py-8 max-w-4xl">
      <div className="space-y-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-3xl font-bold text-neutral-700 mb-2">File Upload Test</h1>
          <p className="text-neutral-500">
            Test the Google Drive file picker integration with various upload button configurations.
          </p>
        </div>

        {/* Upload Button Examples */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Basic Upload Button */}
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h2 className="text-xl font-semibold mb-4 text-gray-600">Basic Upload Button</h2>
            <p className="text-gray-600 mb-4 text-sm">
              Basic file upload with custom configuration.
            </p>
            <UploadButton
              config={{
                subjectId: selectedSubjectId,
                allowMultiple: true,
                maxFiles: 5,
              }}
              onUploadComplete={handleUploadComplete}
              onUploadStart={handleUploadStart}
              onProgress={handleProgress}
              buttonText="Choose Files from Drive"
              variant="primary"
              size="md"
            />
          </div>

          {/* Single File Upload */}
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h2 className="text-xl font-semibold mb-4 text-gray-600">Single File Upload</h2>
            <p className="text-gray-600 mb-4 text-sm">Pre-configured for single file selection.</p>
            <SingleFileUploadButton
              config={{ subjectId: selectedSubjectId }}
              onUploadComplete={handleUploadComplete}
              onUploadStart={handleUploadStart}
              onProgress={handleProgress}
              variant="secondary"
            />
          </div>

          {/* Multiple Files Upload */}
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h2 className="text-xl font-semibold mb-4 text-gray-600">Multiple Files Upload</h2>
            <p className="text-gray-600 mb-4 text-sm">
              Pre-configured for multiple file selection (up to 10 files).
            </p>
            <MultipleFileUploadButton
              config={{ subjectId: selectedSubjectId }}
              onUploadComplete={handleUploadComplete}
              onUploadStart={handleUploadStart}
              onProgress={handleProgress}
              variant="outline"
            />
          </div>

          {/* Image Upload */}
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h2 className="text-xl font-semibold mb-4 text-gray-600">Image Upload</h2>
            <p className="text-gray-600 mb-4 text-sm">Pre-configured for image files only.</p>
            <ImageUploadButton
              config={{ subjectId: selectedSubjectId }}
              onUploadComplete={handleUploadComplete}
              onUploadStart={handleUploadStart}
              onProgress={handleProgress}
              buttonText="Upload Images"
              variant="primary"
              size="sm"
            />
          </div>

          {/* Document Upload */}
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h2 className="text-xl font-semibold mb-4 text-gray-600">Document Upload</h2>
            <p className="text-gray-600 mb-4 text-sm">
              Pre-configured for documents (PDF, Word, etc.).
            </p>
            <DocumentUploadButton
              config={{ subjectId: selectedSubjectId }}
              onUploadComplete={handleUploadComplete}
              onUploadStart={handleUploadStart}
              onProgress={handleProgress}
              buttonText="Upload Documents"
              variant="ghost"
            />
          </div>
        </div>

        {/* Current Progress */}
        {currentProgress.length > 0 && (
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <h3 className="text-lg font-semibold mb-4">Upload Progress</h3>
            <div className="space-y-2">
              {currentProgress.map((progress, index) => (
                <div key={index} className="flex items-center gap-3 p-3 bg-gray-50 rounded-md">
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-medium">{progress.fileName}</span>
                      <span className="text-xs text-gray-500 capitalize">{progress.status}</span>
                    </div>
                    {progress.status === 'uploading' && (
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${progress.progress}%` }}
                        />
                      </div>
                    )}
                    {progress.error && (
                      <p className="text-xs text-red-600 mt-1">{progress.error}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Upload Results */}
        {uploadResults.length > 0 && (
          <div className="bg-white p-6 rounded-lg shadow-md border">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">Uploaded Files ({uploadResults.length})</h3>
              <button
                onClick={clearResults}
                className="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors"
              >
                Clear
              </button>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {uploadResults.map((file, index) => (
                <div key={index} className="p-4 bg-gray-50 rounded-md border">
                  <div className="flex items-start gap-3">
                    <div className="text-2xl">
                      {file.mime_type.startsWith('image/')
                        ? '🖼️'
                        : file.mime_type.includes('pdf')
                          ? '📄'
                          : file.mime_type.includes('video')
                            ? '🎥'
                            : '📋'}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h4 className="font-medium text-gray-900 truncate">{file.name}</h4>
                      <p className="text-xs text-gray-500 mt-1">Type: {file.mime_type}</p>
                      <p className="text-xs text-gray-500">ID: {file.file_id}</p>
                      <a
                        href={file.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-blue-600 hover:text-blue-800 mt-1 inline-block"
                      >
                        View in Google Drive →
                      </a>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Usage Instructions */}
        <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
          <h3 className="text-lg font-semibold text-blue-900 mb-3">How to Use</h3>
          <div className="text-sm text-blue-800 space-y-2">
            <p>
              <strong>1.</strong> Make sure you have a Google account and have linked it to your PWS
              account.
            </p>
            <p>
              <strong>2.</strong> Click any upload button to open the Google Drive file picker.
            </p>
            <p>
              <strong>3.</strong> Select files from your Google Drive.
            </p>
            <p>
              <strong>4.</strong> The files will be automatically uploaded to your PWS subject.
            </p>
            <p>
              <strong>5.</strong> View the results below to see the uploaded file details.
            </p>
          </div>
        </div>

        {/* Integration Code Examples */}
        <div className="bg-gray-50 p-6 rounded-lg border">
          <h3 className="text-lg font-semibold mb-3">Integration Examples</h3>
          <div className="space-y-4">
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Basic Usage:</h4>
              <pre className="bg-white text-gray-600 p-3 rounded text-xs overflow-x-auto">
                {`<UploadButton
  config={{
    subjectId: "your-subject-id",
    allowMultiple: true,
  }}
  onUploadComplete={(result) => console.log(result)}
/>`}
              </pre>
            </div>

            <div>
              <h4 className="font-medium text-gray-900 mb-2">Pre-configured Buttons:</h4>
              <pre className="bg-white text-gray-600 p-3 rounded text-xs overflow-x-auto">
                {`<SingleFileUploadButton config={{ subjectId: "id" }} />
<MultipleFileUploadButton config={{ subjectId: "id" }} />
<ImageUploadButton config={{ subjectId: "id" }} />
<DocumentUploadButton config={{ subjectId: "id" }} />`}
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
