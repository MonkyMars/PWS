import { useState } from 'react';
import { Link } from 'react-router';
import {
  ArrowLeft,
  Bell,
  FileText,
  Pin,
  Calendar,
  User,
  Download,
  ExternalLink,
  FileImageIcon,
} from 'lucide-react';
import { Button } from '~/components/ui/button';
import { FileViewer } from '~/components/files/file-viewer';
import { useSubject, useAnnouncements, useSubjectFiles } from '~/hooks';
import type { SubjectFile } from '~/types';

interface SubjectDetailProps {
  subjectId: string;
}

export function SubjectDetail({ subjectId }: SubjectDetailProps) {
  const [activeTab, setActiveTab] = useState<'announcements' | 'files'>('announcements');
  const [selectedFile, setSelectedFile] = useState<SubjectFile | null>(null);

  const { data: subject, isLoading: subjectLoading } = useSubject(subjectId);
  // const { data: announcementsData, isLoading: announcementsLoading } = useAnnouncements({
  //   subjectId,
  // });
  const { data: filesData, isLoading: filesLoading } = useSubjectFiles({
    subjectId,
  });

  if (!subject) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-neutral-900 mb-2">Vak niet gevonden</h2>
          <p className="text-neutral-600 mb-4">
            Het opgegeven vak bestaat niet of je hebt geen toegang.
          </p>
          <Link to="/dashboard">
            <Button>Terug naar Dashboard</Button>
          </Link>
        </div>
      </div>
    );
  }

  const getSubjectColor = () => {
    return subject.color || 'var(--color-subject-default)';
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('nl-NL', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'urgent':
        return 'bg-error-100 text-error-800';
      case 'high':
        return 'bg-warning-100 text-warning-800';
      case 'normal':
        return 'bg-primary-100 text-primary-800';
      default:
        return 'bg-neutral-100 text-neutral-800';
    }
  };

  const handleFileClick = (file: SubjectFile) => {
    setSelectedFile(file);
  };

  // const announcements = announcementsData?.items || [];
  const files = filesData?.items || [];

  return (
    <div className="min-h-screen bg-neutral-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/dashboard"
            className="inline-flex items-center text-neutral-600 hover:text-neutral-900 mb-4 transition-colors"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Terug naar Dashboard
          </Link>

          <div className="flex items-center space-x-4 mb-4">
            <div className="w-6 h-6 rounded-full" style={{ backgroundColor: getSubjectColor() }} />
            <div>
              <h1 className="text-3xl font-bold text-neutral-900">{subject.name}</h1>
              <p className="text-neutral-600">
                {subject.code} • {subject.teacherName}
              </p>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="mb-8">
          <div className="border-b border-neutral-200">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab('announcements')}
                className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'announcements'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:border-neutral-300'
                }`}
              >
                <Bell className="h-4 w-4 inline mr-2" />
                Mededelingen
                {/*({announcements.length})*/}
              </button>
              <button
                onClick={() => setActiveTab('files')}
                className={`py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === 'files'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:border-neutral-300'
                }`}
              >
                <FileText className="h-4 w-4 inline mr-2" />
                Bestanden ({files.length})
              </button>
            </nav>
          </div>
        </div>

        {/* Content */}
        {/*{activeTab === 'announcements' && (
          <div className="space-y-6">
            {announcementsLoading ? (
              <div className="flex justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            ) : announcements.length > 0 ? (
              announcements.map((announcement) => (
                <div
                  key={announcement.id}
                  className="bg-white rounded-lg border border-neutral-200 p-6"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center space-x-3">
                      {announcement.isSticky && <Pin className="h-5 w-5 text-warning-500" />}
                      <h3 className="text-lg font-semibold text-neutral-900">
                        {announcement.title}
                      </h3>
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(announcement.priority)}`}
                      >
                        {announcement.priority === 'urgent' && 'Urgent'}
                        {announcement.priority === 'high' && 'Hoog'}
                        {announcement.priority === 'normal' && 'Normaal'}
                        {announcement.priority === 'low' && 'Laag'}
                      </span>
                    </div>
                  </div>

                  <div className="prose prose-sm max-w-none text-neutral-700 mb-4">
                    {announcement.content}
                  </div>

                  <div className="flex items-center justify-between text-sm text-neutral-500">
                    <div className="flex items-center space-x-1">
                      <User className="h-4 w-4" />
                      <span>{announcement.authorName}</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <Calendar className="h-4 w-4" />
                      <span>{formatDate(announcement.createdAt)}</span>
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
                <Bell className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-neutral-900 mb-2">Geen mededelingen</h3>
                <p className="text-neutral-600">
                  Er zijn nog geen mededelingen geplaatst voor dit vak.
                </p>
              </div>
            )}
          </div>
        )}*/}

        {activeTab === 'files' && (
          <div className="space-y-6">
            {filesLoading ? (
              <div className="flex justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            ) : files.length > 0 ? (
              <div className="grid grid-cols-1">
                {files.map((file) => {
                  const isImage = file.mimeType.startsWith('image/');
                  const isPdf = file.mimeType === 'application/pdf';
                  const isText = file.mimeType.startsWith('text/');

                  const className: string = 'w-8 h-8 text-secondary-500';
                  let icon: React.ReactNode;
                  if (isImage) {
                    icon = <FileImageIcon className={className} />;
                  } else if (isPdf) {
                    icon = <FileText className={className} />;
                  } else if (isText) {
                    icon = <FileText className={className} />;
                  } else {
                    icon = <FileText className={className} />;
                  }

                  return (
                    <div
                      key={file.id}
                      className="flex items-center justify-between py-4 px-6 border-b border-neutral-200 hover:bg-neutral-50 transition-colors cursor-pointer"
                      onClick={() => handleFileClick(file)}
                    >
                      <div className="flex items-center space-x-4 flex-1 min-w-0">
                        <div className="flex-shrink-0 ">{icon}</div>
                        <div className="flex-1 min-w-0">
                          <h4 className="font-medium text-neutral-900 truncate">{file.name}</h4>
                          <div className="flex items-center space-x-4 mt-1 text-sm text-neutral-500">
                            {file.createdAt && (
                              <div className="flex items-center space-x-1">
                                <Calendar className="h-3 w-3" />
                                <span>{formatDate(file.createdAt)}</span>
                              </div>
                            )}
                            {file.uploaderId && (
                              <div className="flex items-center space-x-1">
                                <User className="h-3 w-3" />
                                <span>{file.uploaderId}</span>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center space-x-2 flex-shrink-0">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleFileClick(file);
                          }}
                        >
                          <ExternalLink className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="bg-white rounded-lg border border-neutral-200 p-12 text-center">
                <FileText className="h-16 w-16 text-neutral-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-neutral-900 mb-2">Geen bestanden</h3>
                <p className="text-neutral-600">
                  Er zijn nog geen bestanden geüpload voor dit vak.
                </p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* File Viewer Modal */}
      {selectedFile && (
        <FileViewer
          file={selectedFile}
          isOpen={!!selectedFile}
          onClose={() => setSelectedFile(null)}
        />
      )}
    </div>
  );
}
