create table if not exists public.deadlines (
	id uuid not null default gen_random_uuid (),
	subject_id uuid not null,
	owner_id uuid not null,
	title text not null,
	description text null,
	due_date timestamp with time zone not null,
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now(),
	constraint deadlines_pkey primary key (id),
	constraint fk_deadlines_subjects foreign key (subject_id) references public.subjects (id) on delete cascade,
	constraint fk_deadlines_users foreign key (owner_id) references public.users (id) on delete cascade
) TABLESPACE pg_default;

create index IF not exists idx_deadlines_owner_id on public.deadlines using btree (owner_id) TABLESPACE pg_default;

create index IF not exists idx_deadlines_due_date on public.deadlines using btree (due_date) TABLESPACE pg_default;

create index IF not exists idx_deadlines_subject_id on public.deadlines using btree (subject_id) TABLESPACE pg_default;

create index IF not exists idx_deadlines_created_at on public.deadlines using btree (created_at) TABLESPACE pg_default;

create index IF not exists idx_deadlines_updated_at on public.deadlines using btree (updated_at) TABLESPACE pg_default;
