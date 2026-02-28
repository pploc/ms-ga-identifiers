-- Create login_attempts table
CREATE TABLE login_attempts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identity_id  UUID REFERENCES identities(id) ON DELETE SET NULL,
    email        VARCHAR(255) NOT NULL,
    ip_address   VARCHAR(45),
    success      BOOLEAN NOT NULL,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_login_attempts_identity_id ON login_attempts(identity_id);
CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_attempted_at ON login_attempts(attempted_at);
