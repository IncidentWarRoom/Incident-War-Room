-- =====================
-- EVENT MEDIA
-- =====================

ALTER TABLE incident_events
    ADD COLUMN IF NOT EXISTS media_url TEXT;
