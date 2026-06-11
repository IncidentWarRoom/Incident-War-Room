CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================
-- INCIDENTS
-- =====================

CREATE TYPE incident_status AS ENUM (
    'active',
    'investigating',
    'mitigated',
    'resolved'
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
    status      incident_status   NOT NULL DEFAULT 'active',
    chat_id     BIGINT NOT NULL,
    created_by  BIGINT NOT NULL,                -- tg_user_id
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_incidents_one_active_per_chat
    ON incidents (chat_id)
    WHERE status != 'resolved';

CREATE INDEX idx_incidents_chat_id
    ON incidents (chat_id, created_at DESC);

CREATE INDEX idx_incidents_status
    ON incidents (status)
    WHERE status != 'resolved';

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
    author_id   BIGINT,                         
    username    TEXT,                           
    message     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_incident_timeline
    ON incident_events (incident_id, created_at ASC);

CREATE INDEX idx_events_incident_type
    ON incident_events (incident_id, type);