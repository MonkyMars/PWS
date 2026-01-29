import { useState, useEffect } from 'react';
import { Link } from 'react-router';
import {
  ArrowLeft,
  Bell,
  FileText,
  Calendar,
  User,
  ExternalLink,
  FileImageIcon,
  Folder,
  ChevronLeft,
  ChevronRight,
  Keyboard,
  HelpCircle,
  ChevronDown,
  ChevronUp,
  Printer,
  ArrowRight,
} from 'lucide-react';
import { Button } from '~/components/ui/button';
import { FileViewer } from '~/components/files/file-viewer';
import { useSubject, useSubjectFiles, useSubjectFolders, useSubjectTeachers } from '~/hooks';
import type { SubjectFile } from '~/types';
import { useDeadlines } from '~/hooks/use-deadlines';

import { useAuth } from '~/hooks';

interface SubjectDetailProps {
  subjectId: string;
}

export function SubjectDetail({ subjectId }: SubjectDetailProps) {
  const [activeTab, setActiveTab] = useState<'announcements' | 'files'>('announcements');
  const [selectedFile, setSelectedFile] = useState<SubjectFile | null>(null);
  const [selectedFolder, setSelectedFolder] = useState<string>(subjectId);
  const [folderHistory, setFolderHistory] = useState<string[]>([subjectId]);
  const [folderNames, setFolderNames] = useState<{ [key: string]: string }>({});
  const { user } = useAuth();
  const [showKeyboardHelp, setShowKeyboardHelp] = useState<boolean>(false);
  const { data: subject } = useSubject(subjectId);
  const { data: teachers } = useSubjectTeachers(subjectId);
  const { data: deadlines } = useDeadlines({ subjectId });
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

  const handlePrint = async (item: SubjectFile) => {
    const file = await fetch(`https://drive.usercontent.google.com/download?id=${item.fileId}`);
    const blob = await file.blob();
    const url = URL.createObjectURL(blob);
    const printWindow = window.open(url, '_blank');
    if (printWindow) {
      printWindow.focus();
      printWindow.print();
    } else {
      alert('Pop-up geblokkeerd. Sta pop-ups toe om te kunnen printen.');
    }
  };

  // const announcements = announcementsData?.items || [];
  return (
    <div className="min-h-screen bg-neutral-50">
      <div className="max-w-480 mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
              <p className="text-neutral-600">{subject.code}</p>
            </div>
          </div>
        </div>

        {/* Main Content Grid: 3:1 layout */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Left side - Main content (3/4) */}
          <div className="lg:col-span-3">
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
                        <div className="flex items-center">
                          <button
                            onClick={() => setShowKeyboardHelp(!showKeyboardHelp)}
                            className="flex items-center gap-1.5 px-2 py-1 text-xs text-neutral-500 hover:text-neutral-700 hover:bg-neutral-50 rounded transition-colors"
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
                  <div className="bg-neutral-50 border border-neutral-200 rounded-lg p-3 animate-slide-down">
                    <div className="flex items-start gap-2.5">
                      <HelpCircle className="h-4 w-4 text-neutral-400 mt-0.5 shrink-0" />
                      <div className="flex-1">
                        <h4 className="font-medium text-neutral-700 mb-2 text-sm">Sneltoetsen</h4>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
                          <div className="space-y-1.5">
                            <div className="flex items-center gap-2">
                              <kbd className="px-1.5 py-0.5 bg-white border border-neutral-300 text-neutral-600 rounded text-xs font-mono min-w-[1.5rem] text-center">
                                1-9
                              </kbd>
                              <span className="text-neutral-600">
                                Open map/bestand (in volgorde)
                              </span>
                            </div>
                            <div className="flex items-center gap-2">
                              <kbd className="px-1.5 py-0.5 bg-white border border-neutral-300 text-neutral-600 rounded text-xs font-mono min-w-[1.5rem] text-center">
                                ESC
                              </kbd>
                              <span className="text-neutral-600">Ga terug naar vorige map</span>
                            </div>
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
                          className="group flex items-center justify-between bg-white rounded-lg p-2 border-b border-neutral-200 hover:bg-neutral-50 transition-colors cursor-pointer relative"
                          onClick={handleItemClick}
                        >
                          <div className="flex items-center space-x-4 flex-1 min-w-0">
                            <div className="shrink-0">{icon}</div>
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
                          <div className="flex items-center gap-3 shrink-0">
                            {!isFolder && (
                              <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                                <button
                                  title="Printen"
                                  className="p-1 rounded hover:bg-neutral-100 text-neutral-500 hover:text-neutral-700 transition-colors"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handlePrint(item);
                                  }}
                                >
                                  <Printer className="h-3.5 w-3.5" />
                                </button>
                                <button
                                  title="Openen"
                                  className="p-1 rounded hover:bg-neutral-100 text-neutral-500 hover:text-neutral-700 transition-colors"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleFileClick(item);
                                  }}
                                >
                                  <ExternalLink className="h-3.5 w-3.5" />
                                </button>
                              </div>
                            )}
                            {showShortcut && (
                              <div className="ml-1 w-5 h-5 flex items-center justify-center bg-neutral-50 border border-neutral-200 rounded text-xs font-mono text-neutral-500">
                                <span className="sr-only">Sneltoets {shortcutNumber}</span>
                                <span>{shortcutNumber}</span>
                              </div>
                            )}
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

          {/* Right side - Info Panel (1/4) */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg border border-neutral-200 p-6 sticky top-8">
              <h3 className="text-lg font-semibold text-neutral-900 mb-4">Vak Informatie</h3>

              {/* Subject Color */}
              <div className="mb-4">
                <p className="text-sm font-medium text-neutral-600 mb-2">Kleur</p>
                <div className="flex items-center space-x-2">
                  <div
                    className="w-8 h-8 rounded-full border-2 border-neutral-200"
                    style={{ backgroundColor: getSubjectColor() }}
                  />
                  <span className="text-sm text-neutral-700">{subject.code}</span>
                </div>
              </div>

              {/* Teachers */}
              {teachers && teachers.length > 0 && (
                <div className="mb-4">
                  <p className="text-sm font-medium text-neutral-600 mb-2">Docenten</p>
                  <div className="space-y-2">
                    {teachers.map((teacher) => (
                      <div
                        key={teacher.id}
                        className="flex items-center space-x-2 text-sm text-neutral-700"
                      >
                        <User className="h-4 w-4 text-neutral-400" />
                        <span>{teacher.username}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Subject Code */}
              <div className="mb-4">
                <p className="text-sm font-medium text-neutral-600 mb-2">Vakcode</p>
                <p className="text-sm text-neutral-700">{subject.code}</p>
              </div>

              {/* Quick Stats */}
              <div className="pt-4 border-t border-neutral-200 space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-neutral-600">Bestanden</span>
                  <span className="font-medium text-neutral-900">{files.length}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-neutral-600">Mappen</span>
                  <span className="font-medium text-neutral-900">{folders.length}</span>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg border border-neutral-200 p-6 sticky mt-8">
              <h3 className="text-lg font-semibold text-neutral-900 mb-4">Deadlines</h3>
              <div className="space-y-2">
                {deadlines && deadlines.length > 0 ? (
                  deadlines.slice(0, 5).map((deadline) => (
                    <div
                      key={deadline.id}
                      className="flex items-center justify-between bg-neutral-50 rounded-lg px-3 py-2 border border-neutral-100 hover:border-primary-200 transition-colors"
                    >
                      <div>
                        <p className="font-medium text-neutral-900 truncate">{deadline.title}</p>
                        <p className="text-xs text-neutral-500">{formatDate(deadline.dueDate)}</p>
                      </div>
                      {/* Role-based hand-in link */}
                      {user?.role === 'student' ? (
                        <Link
                          to={`/deadlines/${deadline.id}/submission`}
                          className="ml-3 flex items-center gap-1 px-2 py-1 rounded text-primary-600 hover:text-white hover:bg-primary-600 text-xs font-medium transition-colors"
                          title="Naar inleverpagina"
                        >
                          <span>Inleveren</span>
                          <ArrowRight className="h-4 w-4" />
                        </Link>
                      ) : (
                        <Link
                          to={`/deadlines/${deadline.id}/submissions`}
                          className="ml-3 flex items-center gap-1 px-2 py-1 rounded text-primary-600 hover:text-white hover:bg-primary-600 text-xs font-medium transition-colors"
                          title="Bekijk ingeleverde opdrachten"
                        >
                          <span>Bekijk ingeleverd</span>
                          <ArrowRight className="h-4 w-4" />
                        </Link>
                      )}
                    </div>
                  ))
                ) : (
                  <div className="text-sm text-neutral-500 text-center py-4">
                    Geen aankomende deadlines.
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
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
