CREATE TABLE IF NOT EXISTS meetings(
    id BIGSERIAL PRIMARY KEY,
    creator_id BIGINT NOT NULL,
    creator_name TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,

    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    meeting_time TIMESTAMPTZ NOT NULL,
    max_people INTEGER
);

CREATE TABLE IF NOT EXISTS meeting_participants (
    meeting_id BIGINT NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    user_name TEXT NOT NULL,

    PRIMARY KEY (meeting_id, user_id)
);