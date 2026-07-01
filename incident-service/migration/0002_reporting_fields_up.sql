-- =====================
-- INCIDENT REPORTING FIELDS
-- =====================

ALTER TABLE incidents
    ADD COLUMN IF NOT EXISTS telegraph_urls JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS report_url     TEXT;
