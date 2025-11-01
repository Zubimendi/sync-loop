-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK (role IN ('owner','admin','member')) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workspace (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    plan TEXT CHECK (plan IN ('free','starter','growth','enterprise')) NOT NULL,
    stripe_customer_id TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS connector (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type TEXT CHECK (type IN ('pg','mysql','s3','excel','sf','gsheets','rest')) NOT NULL,
    config_json JSONB NOT NULL,
    created_by_user_id UUID REFERENCES users(id),
    workspace_id UUID REFERENCES workspace(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sync_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connector_id UUID REFERENCES connector(id) ON DELETE CASCADE,
    schedule_cron TEXT,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ,
    status TEXT CHECK (status IN ('active','paused','error')) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sync_run (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES sync_job(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    status TEXT CHECK (status IN ('running','success','failed','cancelled')) NOT NULL,
    rows_read BIGINT,
    rows_written BIGINT,
    checksum VARCHAR(64),
    log_url TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS destination (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connector_id UUID REFERENCES connector(id) ON DELETE CASCADE,
    type TEXT CHECK (type IN ('pg','s3','excel','gsheets','bq')) NOT NULL,
    config_json JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS field_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES sync_job(id) ON DELETE CASCADE,
    source_field TEXT NOT NULL,
    dest_field TEXT NOT NULL,
    transform TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workspace_user (
    workspace_id UUID REFERENCES workspace(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role TEXT CHECK (role IN ('owner','admin','member')) NOT NULL,
    PRIMARY KEY (workspace_id, user_id)
);

CREATE INDEX idx_sync_run_job_id ON sync_run(job_id);
CREATE INDEX idx_connector_workspace_id ON connector(workspace_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS field_mapping;
DROP TABLE IF EXISTS destination;
DROP TABLE IF EXISTS sync_run;
DROP TABLE IF EXISTS sync_job;
DROP TABLE IF EXISTS connector;
DROP TABLE IF EXISTS workspace_user;
DROP TABLE IF EXISTS workspace;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd