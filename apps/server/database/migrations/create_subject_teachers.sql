create table public.subject_teachers (
  id uuid not null default gen_random_uuid (),
  subject_id uuid not null,
  user_id uuid not null,
  created_at timestamp with time zone not null default now(),
  constraint subject_teachers_pkey primary key (id),
  constraint subject_teachers_subject_id_fkey foreign key (subject_id) references subjects (id) on delete cascade,
  constraint subject_teachers_user_id_fkey foreign key (user_id) references users (id) on delete cascade,
  constraint subject_teachers_unique unique (subject_id, user_id)
) tablespace pg_default;

-- Create index for faster lookups
create index subject_teachers_subject_id_idx on public.subject_teachers using btree (subject_id);
create index subject_teachers_user_id_idx on public.subject_teachers using btree (user_id);
