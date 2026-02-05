import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useParams, useNavigate } from 'react-router';
import { useAuth } from '~/hooks';
import { useDeadlines } from '~/hooks/use-deadlines';
import { useStudentSubmission } from '~/hooks/use-student-submission';
import { apiClient } from '~/lib/api-client';
import { Button } from '~/components/ui/button';
import { Input } from '~/components/ui/input';
import { AlertCircle, FileText, UploadCloud, Loader2, CheckCircle2 } from 'lucide-react';

const MAX_FILES = 10;

export default function StudentSubmissionPage() {
  const { id: deadlineId } = useParams<{ id: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();

  // Redirect teachers/admins to the submissions list for this deadline
  useEffect(() => {
    if (user && (user.role === 'teacher' || user.role === 'admin')) {
      navigate(`/deadlines/${deadlineId}/submissions`, { replace: true });
    }
  }, [user, deadlineId, navigate]);

  // Fetch deadline details
  const { data: deadlines = [], isLoading: deadlinesLoading } = useDeadlines();
  const deadline = deadlines.find((d) => d.id === deadlineId);

  // Fetch student's current submission
  const {
    data: submission,
    isLoading: submissionLoading,
    refetch: refetchSubmission,
  } = useStudentSubmission(deadlineId);

  // Form state
  const [message, setMessage] = useState('');
  const [files, setFiles] = useState<File[]>([]);
  const [filePreviews, setFilePreviews] = useState<string[]>([]);
  const [uploading, setUploading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Populate form with existing submission
  useEffect(() => {
    if (submission) {
      setMessage(submission.message || '');
      // Note: We can't fetch actual files from Google Drive by ID here, so just show count
      setFilePreviews(submission.fileIds || []);
    }
  }, [submission]);

  // File input ref for manual trigger
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Handle file selection
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) return;
    const selected = Array.from(e.target.files);
    if (selected.length + files.length > MAX_FILES) {
      setError(`Je mag maximaal ${MAX_FILES} bestanden uploaden.`);
      return;
    }
    setFiles((prev) => [...prev, ...selected]);
    setError(null);
  };

  // Handle drag-and-drop
  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    const dropped = Array.from(e.dataTransfer.files);
    if (dropped.length + files.length > MAX_FILES) {
      setError(`Je mag maximaal ${MAX_FILES} bestanden uploaden.`);
      return;
    }
    setFiles((prev) => [...prev, ...dropped]);
    setError(null);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  // Remove a file from the list
  const handleRemoveFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index));
  };

  // Submit handler
  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      setUploading(true);
      setError(null);
      setSuccess(false);

      try {
        // 1. Upload files to backend (Google Drive integration)
        let fileIds: string[] = submission?.fileIds ? [...submission.fileIds] : [];
        if (files.length > 0) {
          for (const file of files) {
            // Use your backend's file upload endpoint (should return Google Drive file ID)
            const response = await apiClient.uploadFile<{ fileId: string }>('/files/upload', file, {
              deadlineId,
            });
            if (!response.success || !response.data?.fileId) {
              throw new Error(response.message || 'Fout bij uploaden bestand');
            }
            fileIds.push(response.data.fileId);
          }
        }

        // 2. Submit or update the submission
        const res = await apiClient.post(`/deadlines/${deadlineId}/submission`, {
          file_ids: fileIds,
          message,
        });

        if (!res.success) {
          throw new Error(res.message || 'Fout bij inleveren opdracht');
        }

        setSuccess(true);
        setFiles([]);
        refetchSubmission();
      } catch (err: any) {
        setError(err.message || 'Onbekende fout');
      } finally {
        setUploading(false);
      }
    },
    [files, message, deadlineId, submission, refetchSubmission]
  );

  // Status indicator
  let statusNode: React.ReactNode = null;
  if (submissionLoading) {
    statusNode = (
      <span className="text-xs px-2 py-1 bg-neutral-100 text-neutral-700 rounded">Laden...</span>
    );
  } else if (submission) {
    if (submission.isLate) {
      statusNode = (
        <span className="text-xs px-2 py-1 bg-red-100 text-red-700 rounded">Te laat</span>
      );
    } else {
      statusNode = (
        <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded">
          Op tijd ingeleverd
        </span>
      );
    }
    if (submission.isUpdated) {
      statusNode = (
        <>
          {statusNode}
          <span className="ml-2 text-xs px-2 py-1 bg-yellow-100 text-yellow-800 rounded">
            Gewijzigd na deadline
          </span>
        </>
      );
    }
  } else {
    statusNode = (
      <span className="text-xs px-2 py-1 bg-neutral-100 text-neutral-700 rounded">
        Nog niet ingeleverd
      </span>
    );
  }

  if (!user || user.role !== 'student') {
    // The redirect above will handle teachers/admins, but show nothing here for non-students
    return null;
  }

  if (deadlinesLoading) {
    return <div className="p-8 text-center text-neutral-500">Laden...</div>;
  }

  if (!deadline) {
    return (
      <div className="p-8 text-center text-red-600">
        <AlertCircle className="inline-block mr-2" />
        Opdracht niet gevonden.
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto py-10 px-4">
      <div className="mb-8">
        <h1 className="text-2xl font-bold mb-2">{deadline.title}</h1>
        <div className="text-neutral-600 mb-2">{deadline.description}</div>
        <div className="text-sm text-neutral-500 mb-2">
          Deadline: {new Date(deadline.dueDate).toLocaleString('nl-NL')}
        </div>
        <div className="mb-2">{statusNode}</div>
        <Button variant="outline" onClick={() => navigate('/deadlines/mine')}>
          Terug naar opdrachten
        </Button>
      </div>

      <form
        onSubmit={handleSubmit}
        className="bg-white rounded-lg border border-neutral-200 p-6 space-y-6"
      >
        <div>
          <label className="block font-medium mb-1">Bericht aan docent (optioneel)</label>
          <Input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Voeg een bericht toe aan je inlevering..."
            disabled={uploading}
          />
        </div>

        <div>
          <label className="block font-medium mb-1">Bestanden uploaden (max {MAX_FILES})</label>
          <div
            className="border-2 border-dashed border-neutral-300 rounded-lg p-4 mb-2 flex flex-col items-center justify-center cursor-pointer hover:border-primary-400 transition-colors"
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onClick={() => fileInputRef.current?.click()}
            style={{ minHeight: 80 }}
          >
            <UploadCloud className="h-8 w-8 text-neutral-400 mb-2" />
            <div className="text-neutral-500 mb-1">
              Sleep bestanden hierheen of klik om te selecteren
            </div>
            <div className="text-xs text-neutral-400">
              Ondersteund: docx, Google Docs, afbeeldingen, etc.
            </div>
            <input
              ref={fileInputRef}
              type="file"
              multiple
              className="hidden"
              onChange={handleFileChange}
              disabled={uploading}
            />
          </div>
          <div className="space-y-1">
            {files.map((file, idx) => (
              <div
                key={idx}
                className="flex items-center justify-between bg-neutral-100 rounded px-3 py-1"
              >
                <span className="text-sm">{file.name}</span>
                <Button
                  type="button"
                  size="sm"
                  variant="ghost"
                  onClick={() => handleRemoveFile(idx)}
                  disabled={uploading}
                >
                  Verwijder
                </Button>
              </div>
            ))}
            {/* Show previously uploaded files (from submission) */}
            {filePreviews.length > 0 && (
              <div className="mt-2">
                <div className="text-xs text-neutral-500 mb-1">Eerder ingeleverde bestanden:</div>
                {filePreviews.map((fid, idx) => (
                  <div
                    key={fid}
                    className="flex items-center gap-2 text-xs bg-neutral-50 rounded px-2 py-1 mb-1"
                  >
                    <FileText className="h-4 w-4 text-neutral-400" />
                    <span>
                      Bestand {idx + 1} (Google Drive ID: {fid})
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {error && (
          <div className="text-red-600 flex items-center gap-2">
            <AlertCircle className="h-4 w-4" /> {error}
          </div>
        )}
        {success && (
          <div className="text-green-700 flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4" /> Inlevering succesvol opgeslagen!
          </div>
        )}

        <Button
          type="submit"
          variant="primary"
          className="w-full flex items-center justify-center"
          disabled={uploading}
        >
          {uploading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Bezig met uploaden...
            </>
          ) : submission ? (
            'Bijwerken'
          ) : (
            'Inleveren'
          )}
        </Button>
      </form>
    </div>
  );
}
