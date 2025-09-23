create table public.subjects (
  id uuid not null default gen_random_uuid (),
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),
  name text not null,
  constraint subjects_pkey primary key (id)
) TABLESPACE pg_default;
