CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,

    -- Store a BCrypt/Argon2 hash, NEVER the plain password
    password_hash VARCHAR(255) NOT NULL,

    name VARCHAR(255) NOT NULL,

    created_requests_count BIGINT NOT NULL DEFAULT 0,
    completed_requests_count BIGINT NOT NULL DEFAULT 0
);
