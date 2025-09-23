create table public.users (
  id uuid not null default gen_random_uuid (),
  created_at timestamp with time zone not null default now(),
  username text null,
  email text null,
  role text null,
  password_hash text null,
  constraint users_pkey primary key (id)
) TABLESPACE pg_default;
