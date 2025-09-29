# File Upload Components

This directory contains modular file upload components that integrate with Google Drive and your PWS API.

## Overview

The upload system works by:
1. Opening a Google Drive file picker
2. Allowing users to select files from their Google Drive
3. Sending the selected file metadata to your PWS API
4. Storing file references in your database

## Components

### `UploadButton`

The main upload button component with full customization options.

```tsx
import { UploadButton } from '~/components/files/upload-button';

<UploadButton
  config={{
    subjectId: "your-subject-id",
    allowMultiple: true,
    maxFiles: 5,
    acceptedMimeTypes: ['application/pdf', 'image/*']
  }}
  onUploadComplete={(result) => console.log(result)}
  onUploadStart={() => console.log('Upload started')}
  onProgress={(progress) => console.log(progress)}
  buttonText="Upload Files"
  variant="primary"
  size="md"
/>
```

### Pre-configured Components

These components come with sensible defaults for common use cases:

#### `SingleFileUploadButton`
```tsx
import { SingleFileUploadButton } from '~/components/files/upload-button';

<SingleFileUploadButton
  config={{ subjectId: "your-subject-id" }}
  onUploadComplete={handleUpload}
/>
```

#### `MultipleFileUploadButton`
```tsx
import { MultipleFileUploadButton } from '~/components/files/upload-button';

<MultipleFileUploadButton
  config={{ subjectId: "your-subject-id" }}
  onUploadComplete={handleUpload}
/>
```

#### `ImageUploadButton`
```tsx
import { ImageUploadButton } from '~/components/files/upload-button';

<ImageUploadButton
  config={{ subjectId: "your-subject-id" }}
  onUploadComplete={handleUpload}
/>
```

#### `DocumentUploadButton`
```tsx
import { DocumentUploadButton } from '~/components/files/upload-button';

<DocumentUploadButton
  config={{ subjectId: "your-subject-id" }}
  onUploadComplete={handleUpload}
/>
```

## Configuration

### `UploadConfig`

```tsx
interface UploadConfig {
  subjectId: string;           // Required: Subject ID for file association
  allowMultiple?: boolean;     // Allow multiple file selection
  maxFiles?: number;          // Maximum number of files
  acceptedMimeTypes?: string[]; // Allowed MIME types
  maxFileSize?: number;       // Maximum file size in bytes
}
```

### Button Props

```tsx
interface UploadButtonProps {
  config: UploadConfig;
  onUploadComplete?: (result: UploadResult) => void;
  onUploadStart?: () => void;
  onProgress?: (progress: FileUploadProgress[]) => void;
  buttonText?: string;
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
}
```

## Usage Patterns

### Basic Upload

```tsx
function MyComponent() {
  const handleUpload = (result: UploadResult) => {
    if (result.success) {
      console.log('Uploaded files:', result.files);
    } else {
      console.error('Upload failed:', result.errors);
    }
  };

  return (
    <SingleFileUploadButton
      config={{ subjectId: "abc123" }}
      onUploadComplete={handleUpload}
    />
  );
}
```

### With Progress Tracking

```tsx
function MyComponent() {
  const [progress, setProgress] = useState<FileUploadProgress[]>([]);

  return (
    <div>
      <MultipleFileUploadButton
        config={{ subjectId: "abc123" }}
        onProgress={setProgress}
        onUploadComplete={handleUpload}
      />
      
      {progress.map((p, i) => (
        <div key={i}>
          {p.fileName}: {p.status} ({Math.round(p.progress)}%)
        </div>
      ))}
    </div>
  );
}
```

### Custom Configuration

```tsx
function MyComponent() {
  return (
    <UploadButton
      config={{
        subjectId: "abc123",
        allowMultiple: true,
        maxFiles: 3,
        acceptedMimeTypes: ['application/pdf', 'image/jpeg', 'image/png']
      }}
      buttonText="Upload Homework"
      variant="outline"
      size="lg"
      onUploadComplete={handleUpload}
    />
  );
}
```

## Creating Custom Upload Buttons

Use the `createUploadButton` higher-order component to create domain-specific buttons:

```tsx
import { createUploadButton } from '~/components/files/upload-button';

// Create a homework upload button with predefined config
export const HomeworkUploadButton = createUploadButton({
  allowMultiple: true,
  maxFiles: 5,
  acceptedMimeTypes: ['application/pdf', 'image/*']
});

// Use it
<HomeworkUploadButton 
  config={{ subjectId: "abc123" }}
  onUploadComplete={handleUpload}
/>
```

## Swapping Buttons

The modular design allows you to easily swap between different upload button types:

```tsx
function MyComponent({ uploadType }: { uploadType: 'single' | 'multiple' | 'image' }) {
  const commonProps = {
    config: { subjectId: "abc123" },
    onUploadComplete: handleUpload,
  };

  // Just change the component - no other changes needed!
  switch (uploadType) {
    case 'single':
      return <SingleFileUploadButton {...commonProps} />;
    case 'multiple':
      return <MultipleFileUploadButton {...commonProps} />;
    case 'image':
      return <ImageUploadButton {...commonProps} />;
  }
}
```

## Prerequisites

1. **Google Account Linking**: Users must have linked their Google account through the `/auth/google/url` endpoint
2. **Environment Variables**: Set `VITE_GOOGLE_API_KEY` in your environment
3. **Authentication**: User must be logged in (admin or teacher role required for uploads)

## API Integration

The components automatically handle:
- Google OAuth token management
- File selection via Google Drive Picker
- API calls to `/files/upload/single` and `/files/upload/multiple`
- Error handling and retry logic

## File Structure

```
files/
├── upload-button.tsx          # Main upload components
├── upload-example.tsx         # Example implementation
├── file-viewer.tsx           # File viewing component
└── README.md                 # This documentation
```

## Testing

Visit `/test-upload` to see all upload button variations in action and test the functionality.

## Error Handling

The components handle common error scenarios:
- User not authenticated
- Google account not linked
- API errors
- Network failures
- File validation errors

Errors are displayed in the UI and passed to the `onUploadComplete` callback.

## Styling

All components use Tailwind CSS classes and support:
- Custom `className` props
- Variant styles (primary, secondary, outline, ghost)
- Size options (sm, md, lg)
- Custom button text
- Progress indicators
- Error states

## Security

- Uses OAuth 2.0 for Google Drive access
- Files remain in user's Google Drive (only metadata stored in PWS)
- JWT authentication required for API calls
- CSRF protection via state parameters