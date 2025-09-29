create table public.user_subjects (
  id uuid not null default gen_random_uuid (),
  user_id uuid not null,
  subject_id uuid not null,
  constraint user_subjects_pkey primary key (id),
  constraint user_subjects_subject_id_fkey foreign KEY (subject_id) references subjects (id) on delete CASCADE,
  constraint user_subjects_user_id_fkey foreign KEY (user_id) references users (id) on delete CASCADE
) TABLESPACE pg_default;
