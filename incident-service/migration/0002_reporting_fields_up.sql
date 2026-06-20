-- =====================
-- INCIDENT REPORTING FIELDS
-- =====================

ALTER TABLE incidents
    ADD COLUMN telegraph_urls JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN report_url     TEXT;
