/**
 * Subject interface representing academic subjects in the ELO system
 */
export interface Subject {
  id: string;
  name: string;
  description: string;
  code: string; // e.g., "WISK", "NATK", "NEDD"
  color: string; // Hex color for visual identification
  teacherId: string;
  teacherName: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * Announcement interface for subject-specific announcements
 */
export interface Announcement {
  id: string;
  subjectId: string;
  title: string;
  content: string;
  priority: 'low' | 'normal' | 'high' | 'urgent';
  isSticky: boolean; // Pinned announcements
  authorId: string;
  authorName: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * File interface for uploaded documents and resources
 */
export interface SubjectFile {
  id: string;
  subjectId: string;
  name: string;
  originalName: string;
  description?: string;
  mimeType: string;
  size: number; // in bytes
  url: string;
  category: 'presentation' | 'document' | 'assignment' | 'resource' | 'other';
  uploaderId: string;
  uploaderName: string;
  isPublic: boolean;
  downloadCount: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * Subject with related data for dashboard display
 */
export interface SubjectWithDetails extends Subject {
  recentAnnouncements: Announcement[];
  fileCount: number;
  lastActivity: string;
}

/**
 * Subject enrollment interface linking users to subjects
 */
export interface SubjectEnrollment {
  id: string;
  subjectId: string;
  userId: string;
  role: 'student' | 'teacher';
  enrolledAt: string;
}
