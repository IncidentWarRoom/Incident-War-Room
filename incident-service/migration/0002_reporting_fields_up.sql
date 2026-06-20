-- =====================
-- INCIDENT REPORTING FIELDS
-- =====================

ALTER TABLE incidents
    ADD COLUMN topic_id       BIGINT,                          -- linked Telegram Topic
    ADD COLUMN telegraph_urls JSONB NOT NULL DEFAULT '[]'::jsonb, -- timeline Telegraph page URLs
    ADD COLUMN report_url     TEXT;                            -- PDF report URL in object storage
