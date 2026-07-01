-- =====================
-- EVENT IMAGE
-- =====================

ALTER TABLE incident_events
    ADD COLUMN IF NOT EXISTS image_url TEXT;
