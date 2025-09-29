-- Migration: Create audit_logs table
-- Description: Creates the audit_logs table for storing application audit events
-- Author: MonkyMars
-- Date: 2025-09-29

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    level VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    attrs JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_level ON audit_logs(level);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- Create index on JSONB attrs for faster JSON queries
CREATE INDEX IF NOT EXISTS idx_audit_logs_attrs_gin ON audit_logs USING GIN(attrs);

-- Add comments for documentation
COMMENT ON TABLE audit_logs IS 'Application audit log entries for tracking errors and important events';
COMMENT ON COLUMN audit_logs.id IS 'Primary key, auto-incrementing sequence';
COMMENT ON COLUMN audit_logs.timestamp IS 'When the audit event occurred (from application)';
COMMENT ON COLUMN audit_logs.level IS 'Log level (ERROR, WARN, INFO, DEBUG)';
COMMENT ON COLUMN audit_logs.message IS 'Human-readable audit message';
COMMENT ON COLUMN audit_logs.attrs IS 'Additional structured attributes as JSON';
COMMENT ON COLUMN audit_logs.created_at IS 'When the record was inserted into the database';

-- Add constraint to ensure level is valid
ALTER TABLE audit_logs 
ADD CONSTRAINT chk_audit_logs_level 
CHECK (level IN ('ERROR', 'WARN', 'INFO', 'DEBUG'));

-- Add constraint to ensure message is not empty
ALTER TABLE audit_logs 
ADD CONSTRAINT chk_audit_logs_message_not_empty 
CHECK (LENGTH(TRIM(message)) > 0);

-- Add a retention policy comment for future automation
COMMENT ON TABLE audit_logs IS 'Application audit log entries for tracking errors and important events. Retention: 90 days for INFO/DEBUG, 1 year for WARN/ERROR';

-- Set up table permissions (adjust as needed for your user roles)
-- GRANT SELECT, INSERT ON audit_logs TO app_user;
-- GRANT USAGE ON SEQUENCE audit_logs_id_seq TO app_user;

-- Example cleanup query for old logs (run this periodically):
-- DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '90 days' AND level IN ('INFO', 'DEBUG');
-- DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '1 year' AND level IN ('WARN', 'ERROR');