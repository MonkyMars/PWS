create table public.health_logs (
  id uuid not null default gen_random_uuid (),
  timestamp timestamp with time zone not null,
  service text not null,
  status_code integer not null,
  request_count bigint not null default 0,
  error_count bigint not null default 0,
  time_span bigint not null,
  average_latency double precision null,
  constraint health_logs_pkey primary key (id)
) TABLESPACE pg_default;

create index IF not exists idx_health_logs_timestamp on public.health_logs using btree ("timestamp" desc) TABLESPACE pg_default;

create index IF not exists idx_health_logs_service on public.health_logs using btree (service) TABLESPACE pg_default;

create index IF not exists idx_health_logs_service_timestamp on public.health_logs using btree (service, "timestamp" desc) TABLESPACE pg_default;
