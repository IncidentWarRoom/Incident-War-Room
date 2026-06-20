ALTER TABLE incidents
    DROP COLUMN IF EXISTS report_url,
    DROP COLUMN IF EXISTS telegraph_urls;
