-- Create developers table
CREATE TABLE IF NOT EXISTS developers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(100) DEFAULT 'developer',
    avatar VARCHAR(500),
    status VARCHAR(50) DEFAULT 'offline',
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_developers_status ON developers(status);
CREATE INDEX idx_developers_email ON developers(email);
CREATE INDEX idx_developers_last_active ON developers(last_active);
