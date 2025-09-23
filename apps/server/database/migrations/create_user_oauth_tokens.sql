-- Create user_oauth_tokens table for storing OAuth refresh tokens
-- This table stores refresh tokens for OAuth providers (Google, etc.) linked to user accounts

CREATE TABLE IF NOT EXISTS user_oauth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL DEFAULT 'google',
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Ensure one token per user per provider
    UNIQUE(user_id, provider)
);

-- Create index for faster lookups by user_id
CREATE INDEX IF NOT EXISTS idx_user_oauth_tokens_user_id ON user_oauth_tokens(user_id);

-- Create index for faster lookups by provider
CREATE INDEX IF NOT EXISTS idx_user_oauth_tokens_provider ON user_oauth_tokens(provider);

-- Enable Row Level Security
ALTER TABLE user_oauth_tokens ENABLE ROW LEVEL SECURITY;

-- Create policy to allow users to only access their own tokens
CREATE POLICY "Users can only access their own OAuth tokens" ON user_oauth_tokens
    FOR ALL USING (auth.uid() = user_id);

-- Update the updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_user_oauth_tokens_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_oauth_tokens_updated_at
    BEFORE UPDATE ON user_oauth_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_user_oauth_tokens_updated_at();

-- Add comments for documentation
COMMENT ON TABLE user_oauth_tokens IS 'Stores OAuth refresh tokens for external providers linked to user accounts';
COMMENT ON COLUMN user_oauth_tokens.user_id IS 'Reference to the user who owns this OAuth token';
COMMENT ON COLUMN user_oauth_tokens.provider IS 'OAuth provider name (google, microsoft, etc.)';
COMMENT ON COLUMN user_oauth_tokens.refresh_token IS 'Encrypted refresh token for the OAuth provider';
