create table public.files (
  id uuid not null default gen_random_uuid (),
  created_at timestamp with time zone not null default now(),
  file_id text not null,
  name text not null,
  mime_type text not null,
  subject_id uuid not null,
  uploaded_by uuid not null,
  constraint files_pkey primary key (id),
  constraint files_subject_id_fkey foreign KEY (subject_id) references subjects (id) on delete CASCADE,
  constraint files_uploaded_by_fkey foreign KEY (uploaded_by) references users (id)
) TABLESPACE pg_default;
