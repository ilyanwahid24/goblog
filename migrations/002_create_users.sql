CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    token VARCHAR(64) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Insert default admin user: username=admin, password=password
-- $2a$10$wT.fB/M95N3p1/Kz.C5.N... is bcrypt hash for 'password'
INSERT INTO users (username, password_hash)
VALUES ('admin', '$2a$10$f/WfzCtmx5tcQZG1oWhw1OQTwUNXDcThxIMSHKMoL7bSnQGpKwrsi')
ON CONFLICT (username) DO NOTHING;
