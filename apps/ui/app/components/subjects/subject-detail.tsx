import { useState, useEffect } from 'react';
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
  Folder,
  ChevronLeft,
  ChevronRight,
  Keyboard,
  HelpCircle,
  ChevronDown,
  ChevronUp,
} from 'lucide-react';
import { Button } from '~/components/ui/button';
import { FileViewer } from '~/components/files/file-viewer';
import { useSubject, useAnnouncements, useSubjectFiles, useSubjectFolders } from '~/hooks';
import type { SubjectFile } from '~/types';

interface SubjectDetailProps {
  subjectId: string;
}

export function SubjectDetail({ subjectId }: SubjectDetailProps) {
  const [activeTab, setActiveTab] = useState<'announcements' | 'files'>('announcements');
  const [selectedFile, setSelectedFile] = useState<SubjectFile | null>(null);
  const [selectedFolder, setSelectedFolder] = useState<string>(subjectId);
  const [folderHistory, setFolderHistory] = useState<string[]>([subjectId]);
  const [folderNames, setFolderNames] = useState<{ [key: string]: string }>({});
  const [showKeyboardHelp, setShowKeyboardHelp] = useState<boolean>(false);
  const { data: subject, isLoading: subjectLoading } = useSubject(subjectId);
  // const { data: announcementsData, isLoading: announcementsLoading } = useAnnouncements({
  //   subjectId,
  // });
  const { data: filesData, isLoading: filesLoading } = useSubjectFiles({
    subjectId,
    folderId: selectedFolder,
  });
  const { data: folderData, isLoading: foldersLoading } = useSubjectFolders({
    subjectId,
    folderId: selectedFolder,
  });

  // Prepare data for hooks
  const files = filesData?.items || [];
  const folders = folderData?.items || [];
  const combinedItems = [
    ...folders.map((folder) => ({ ...folder, type: 'folder' as const })),
    ...files.map((file) => ({ ...file, type: 'file' as const })),
  ];
  const canGoBack = folderHistory.length > 1;

  // Update folder names when subject is loaded
  useEffect(() => {
    if (subject) {
      setFolderNames((prev) => ({ ...prev, [subjectId]: subject.name }));
    }
  }, [subject, subjectId]);

  // Keyboard shortcuts
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      const target = e.target as HTMLElement;
      const isInInput = ['INPUT', 'TEXTAREA', 'SELECT'].includes(target?.tagName);

      // Don't handle shortcuts if user is typing in an input
      if (isInInput) return;

      // ESC to go back
      if (e.key === 'Escape' && canGoBack) {
        e.preventDefault();
        if (folderHistory.length > 1) {
          const newHistory = [...folderHistory];
          newHistory.pop(); // Remove current folder
          const previousFolder = newHistory[newHistory.length - 1];
          setFolderHistory(newHistory);
          setSelectedFolder(previousFolder);
        }
        return;
      }

      // Only handle number shortcuts on files tab
      if (activeTab !== 'files') return;

      // Navigate to folders/files with number keys (1-9)
      const keyNumber = parseInt(e.key);
      if (!isNaN(keyNumber) && keyNumber >= 1 && keyNumber <= 9 && combinedItems.length > 0) {
        const itemIndex = keyNumber - 1;
        if (combinedItems[itemIndex]) {
          e.preventDefault();
          const item = combinedItems[itemIndex];
          if (item.type === 'folder') {
            setFolderHistory((prev) => [...prev, item.id]);
            setSelectedFolder(item.id);
            setFolderNames((prev) => ({ ...prev, [item.id]: item.name }));
          } else {
            setSelectedFile(item);
          }
        }
      }
    };

    document.addEventListener('keydown', handleGlobalKeyDown);
    return () => document.removeEventListener('keydown', handleGlobalKeyDown);
  }, [activeTab, canGoBack, combinedItems, folderHistory]);

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

  // const getPriorityColor = (priority: string) => {
  //   switch (priority) {
  //     case 'urgent':
  //       return 'bg-error-100 text-error-800';
  //     case 'high':
  //       return 'bg-warning-100 text-warning-800';
  //     case 'normal':
  //       return 'bg-primary-100 text-primary-800';
  //     default:
  //       return 'bg-neutral-100 text-neutral-800';
  //   }
  // };

  const handleFileClick = (file: SubjectFile) => {
    setSelectedFile(file);
  };

  const handleFolderClick = (folderId: string, folderName?: string) => {
    setFolderHistory((prev) => [...prev, folderId]);
    setSelectedFolder(folderId);
    if (folderName) {
      setFolderNames((prev) => ({ ...prev, [folderId]: folderName }));
    }
  };

  const handleBackClick = () => {
    if (folderHistory.length > 1) {
      const newHistory = [...folderHistory];
      newHistory.pop(); // Remove current folder
      const previousFolder = newHistory[newHistory.length - 1];
      setFolderHistory(newHistory);
      setSelectedFolder(previousFolder);
    }
  };

  const handleBreadcrumbClick = (folderId: string) => {
    const index = folderHistory.indexOf(folderId);
    if (index !== -1) {
      const newHistory = folderHistory.slice(0, index + 1);
      setFolderHistory(newHistory);
      setSelectedFolder(folderId);
    }
  };

  // const announcements = announcementsData?.items || [];
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
                Bestanden
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
            {/* Navigation header */}
            <div className="bg-white rounded-lg border border-neutral-200 p-4">
              <div className="flex items-center justify-between">
                {/* Breadcrumb navigation */}
                <div className="flex items-center space-x-2 text-sm">
                  <Folder className="h-4 w-4 text-neutral-500" />
                  {folderHistory.map((folderId, index) => (
                    <div key={folderId} className="flex items-center">
                      {index > 0 && <ChevronRight className="h-4 w-4 mx-2 text-neutral-400" />}
                      <button
                        onClick={() => handleBreadcrumbClick(folderId)}
                        className={`px-2 py-1 rounded transition-colors ${
                          index === folderHistory.length - 1
                            ? 'text-neutral-900 font-medium bg-neutral-100'
                            : 'text-primary-600 hover:text-primary-800 hover:bg-primary-50'
                        }`}
                        disabled={index === folderHistory.length - 1}
                      >
                        {folderNames[folderId] ||
                          (folderId === subjectId ? subject?.name : 'Unknown')}
                      </button>
                    </div>
                  ))}
                </div>

                {/* Back button and shortcuts info */}
                <div className="flex items-center space-x-4">
                  {canGoBack && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleBackClick}
                      className="flex items-center space-x-2"
                    >
                      <ChevronLeft className="h-4 w-4" />
                      <span>Terug</span>
                      <span className="text-xs text-neutral-500 ml-2">(ESC)</span>
                    </Button>
                  )}

                  {combinedItems.length > 0 && (
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => setShowKeyboardHelp(!showKeyboardHelp)}
                        className="flex items-center space-x-2 text-xs text-neutral-500 hover:text-neutral-700 transition-colors"
                      >
                        <Keyboard className="h-3 w-3" />
                        <span>Sneltoetsen</span>
                        {showKeyboardHelp ? (
                          <ChevronUp className="h-3 w-3" />
                        ) : (
                          <ChevronDown className="h-3 w-3" />
                        )}
                      </button>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Keyboard shortcuts help */}
            {showKeyboardHelp && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <div className="flex items-start space-x-3">
                  <HelpCircle className="h-5 w-5 text-blue-500 mt-0.5 flex-shrink-0" />
                  <div className="flex-1">
                    <h4 className="font-medium text-blue-800 mb-2">Sneltoetsen</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                      <div className="space-y-2">
                        <div className="flex items-center space-x-2">
                          <kbd className="px-2 py-1 bg-white border text-neutral-700 border-blue-300 rounded text-xs font-mono">
                            1-9
                          </kbd>
                          <span className="text-blue-800">Open map/bestand (in volgorde)</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <kbd className="px-2 py-1 bg-white border text-neutral-700 border-blue-300 rounded text-xs font-mono">
                            ESC
                          </kbd>
                          <span className="text-blue-800">Ga terug naar vorige map</span>
                        </div>
                      </div>
                      <div className="text-xs text-blue-700">
                        <p className="mb-1">
                          • Grijze nummervakjes in rechterhoek tonen sneltoetsen
                        </p>
                        <p>• Sneltoetsen werken alleen op de Bestanden tab</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {filesLoading || foldersLoading ? (
              <div className="flex justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            ) : combinedItems.length > 0 ? (
              <div className="grid grid-cols-1 p-1 w-full overflow-y-auto rounded-lg border-neutral-200">
                {combinedItems.map((item, index: number) => {
                  const isFolder = item.type === 'folder';
                  const className: string = 'w-8 h-8 text-secondary-500';
                  let icon: React.ReactNode;

                  if (isFolder) {
                    icon = <Folder className={className} />;
                  } else {
                    const isImage = item.mimeType?.startsWith('image/');
                    const isPdf = item.mimeType === 'application/pdf';
                    const isText = item.mimeType?.startsWith('text/');

                    if (isImage) {
                      icon = <FileImageIcon className={className} />;
                    } else if (isPdf) {
                      icon = <FileText className={className} />;
                    } else if (isText) {
                      icon = <FileText className={className} />;
                    } else {
                      icon = <FileText className={className} />;
                    }
                  }

                  const handleItemClick = () => {
                    if (isFolder) {
                      handleFolderClick(item.id, item.name);
                    } else {
                      handleFileClick(item);
                    }
                  };

                  const shortcutNumber = index + 1;
                  const showShortcut = shortcutNumber <= 9;

                  return (
                    <div
                      key={`${item.type}-${item.id}`}
                      className="flex items-center justify-between bg-white rounded-lg p-2 border-b border-neutral-200 hover:bg-neutral-50 transition-colors cursor-pointer relative"
                      onClick={handleItemClick}
                    >
                      <div className="flex items-center space-x-4 flex-1 min-w-0">
                        <div className="flex-shrink-0">{icon}</div>
                        <div className="flex-1 min-w-0">
                          <h4 className="font-medium text-neutral-900 truncate">{item.name}</h4>
                          <div className="flex items-center space-x-4 mt-1 text-sm text-neutral-500">
                            {item.createdAt && (
                              <div className="flex items-center space-x-1">
                                <Calendar className="h-3 w-3" />
                                <span>{formatDate(item.createdAt)}</span>
                              </div>
                            )}
                            {item.uploaderId && (
                              <div className="flex items-center space-x-1">
                                <User className="h-3 w-3" />
                                <span>{item.uploaderId}</span>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                      {!isFolder && (
                        <div className="flex mr-6 items-center space-x-2 flex-shrink-0">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleFileClick(item);
                            }}
                          >
                            <ExternalLink className="h-4 w-4" />
                          </Button>
                        </div>
                      )}
                      {showShortcut && (
                        <div className="absolute bottom-2 right-2 z-10 flex items-center justify-center w-5 h-5 bg-neutral-400/90 text-white text-xs rounded font-medium shadow-sm">
                          {shortcutNumber}
                        </div>
                      )}
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
