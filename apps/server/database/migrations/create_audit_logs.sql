-- Migration: Create audit_logs table
-- Description: Creates the audit_logs table for storing application audit events
-- Author: MonkyMars
-- Date: 2025-09-29

create table public.audit_logs (
  timestamp timestamp with time zone not null default now(),
  level character varying(20) not null,
  message text not null,
  attrs jsonb null,
  created_at timestamp with time zone not null default now(),
  id uuid not null default gen_random_uuid (),
  entry_hash character varying(64) null,
  constraint audit_logs_pkey primary key (id),
  constraint chk_audit_logs_entry_hash_not_empty check (
    (
      (entry_hash is null)
      or (
        length(
          TRIM(
            both
            from
              entry_hash
          )
        ) > 0
      )
    )
  ),
  constraint chk_audit_logs_level check (
    (
      (level)::text = any (
        (
          array[
            'ERROR'::character varying,
            'WARN'::character varying,
            'INFO'::character varying,
            'DEBUG'::character varying
          ]
        )::text[]
      )
    )
  ),
  constraint chk_audit_logs_message_not_empty check (
    (
      length(
        TRIM(
          both
          from
            message
        )
      ) > 0
    )
  )
) TABLESPACE pg_default;

create index IF not exists idx_audit_logs_timestamp on public.audit_logs using btree ("timestamp") TABLESPACE pg_default;

create index IF not exists idx_audit_logs_level on public.audit_logs using btree (level) TABLESPACE pg_default;

create index IF not exists idx_audit_logs_created_at on public.audit_logs using btree (created_at) TABLESPACE pg_default;

create index IF not exists idx_audit_logs_attrs_gin on public.audit_logs using gin (attrs) TABLESPACE pg_default;

create unique INDEX IF not exists idx_audit_logs_entry_hash_unique on public.audit_logs using btree (entry_hash) TABLESPACE pg_default
where
  (entry_hash is not null);
