// Re-export all main components for convenient importing
export { Navigation } from './navigation';
export { Footer } from './footer';

// Auth components
export { LoginForm } from './auth/login-form';
export { RegisterForm } from './auth/register-form';
export { ProtectedRoute } from './auth/protected-route';

// Dashboard components
export { Dashboard } from './dashboard/dashboard';
export { SubjectCard } from './dashboard/subject-card';

// UI components
export { Button } from './ui/button';
export { Input } from './ui/input';

// Subject components
export { SubjectDetail } from './subjects/subject-detail';

// File components
export { FileViewer } from './files/file-viewer';
export {
  UploadButton,
  SingleFileUploadButton,
  MultipleFileUploadButton,
  ImageUploadButton,
  DocumentUploadButton,
} from './files/upload-button';
