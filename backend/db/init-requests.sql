CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,

    creator_id BIGINT NOT NULL,
    creator_name TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT NOT NULL,

    requester_address TEXT NOT NULL,

    type TEXT NOT NULL,

    deadline TIMESTAMPTZ
);