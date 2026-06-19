CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================
-- INCIDENTS
-- =====================

CREATE TYPE incident_status AS ENUM (
    'ACTIVE',
    'CLOSED'
);

CREATE TYPE incident_severity AS ENUM (
    'LOW',
    'MEDIUM',
    'HIGH'
);

CREATE TABLE incidents (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT NOT NULL,
    severity    incident_severity NOT NULL DEFAULT 'MEDIUM',
    status      incident_status   NOT NULL DEFAULT 'ACTIVE',
    chat_id     BIGINT NOT NULL,
    topic_id    BIGINT NOT NULL DEFAULT 0,
    created_by  BIGINT,                         -- tg_user_id
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_incidents_one_active_per_topic
    ON incidents (chat_id, topic_id)
    WHERE status = 'ACTIVE';

CREATE INDEX idx_incidents_chat_id
    ON incidents (chat_id, created_at DESC);

-- =====================
-- EVENTS
-- =====================

CREATE TYPE event_type AS ENUM (
    'INCIDENT_CREATED',
    'COMMENT_ADDED',
    'INCIDENT_CLOSED'
);

CREATE TABLE incident_events (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents (id) ON DELETE CASCADE,
    type        event_type NOT NULL,
    user_id     BIGINT,
    username    TEXT NOT NULL DEFAULT '',
    message     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_incident_timeline
    ON incident_events (incident_id, created_at ASC);