CREATE TABLE IF NOT EXISTS public.submissions (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    deadline_id uuid NOT NULL,
    student_id uuid NOT NULL,
    file_ids text[] NOT NULL, -- Google Drive file IDs
    message text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT submissions_pkey PRIMARY KEY (id),
    CONSTRAINT fk_submissions_deadlines FOREIGN KEY (deadline_id) REFERENCES public.deadlines (id) ON DELETE CASCADE,
    CONSTRAINT fk_submissions_students FOREIGN KEY (student_id) REFERENCES public.users (id) ON DELETE CASCADE,
    CONSTRAINT submissions_unique_per_student_per_deadline UNIQUE (deadline_id, student_id)
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_submissions_deadline_id ON public.submissions USING btree (deadline_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_submissions_student_id ON public.submissions USING btree (student_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON public.submissions USING btree (created_at) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_submissions_updated_at ON public.submissions USING btree (updated_at) TABLESPACE pg_default;

-- Automatically update updated_at timestamp on row update
CREATE OR REPLACE FUNCTION update_submissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_submissions_updated_at ON public.submissions;

CREATE TRIGGER trigger_update_submissions_updated_at
    BEFORE UPDATE ON public.submissions
    FOR EACH ROW
    EXECUTE FUNCTION update_submissions_updated_at();

COMMENT ON TABLE public.submissions IS 'Stores student hand-ins for deadlines, including file references, message, and timestamps';
COMMENT ON COLUMN public.submissions.file_ids IS 'Array of Google Drive file IDs associated with the submission';
COMMENT ON COLUMN public.submissions.message IS 'Plain text message submitted by the student';
COMMENT ON COLUMN public.submissions.deadline_id IS 'Reference to the deadline for this submission';
COMMENT ON COLUMN public.submissions.student_id IS 'Reference to the student (user) who made the submission';
