CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,

    creator_id BIGINT NOT NULL,
    creator_name TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL,

    requester_address TEXT NOT NULL,

    type TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('IN PROGRESS', 'PENDING', 'CANCELLED', 'DONE')),

    deadline TIMESTAMPTZ NOT NULL
);

CREATE TABLE request_applications (
    request_id BIGINT NOT NULL REFERENCES requests(id) ON DELETE CASCADE,

    user_id BIGINT NOT NULL,
    user_name TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'ACCEPTED'
        CHECK (status IN ('ACCEPTED', 'CANCELLED', 'DONE')),

    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (request_id, user_id)
);

CREATE INDEX idx_request_applications_request_id
    ON request_applications(request_id);

CREATE INDEX idx_request_applications_user_id
    ON request_applications(user_id);