#!/bin/sh
set -e

export PGPASSWORD="$POSTGRES_PASSWORD"
PSQL="psql -v ON_ERROR_STOP=1 -h ${POSTGRES_HOST:-postgres} -p ${POSTGRES_PORT:-5432} -U $POSTGRES_USER -d $POSTGRES_DB"

$PSQL -c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now());"

for file in /migrations/*_up.sql; do
    version=$(basename "$file" _up.sql)
    applied=$($PSQL -tA -c "SELECT 1 FROM schema_migrations WHERE version = '$version';")
    if [ "$applied" = "1" ]; then
        echo "migrate: skip $version (already applied)"
        continue
    fi
    echo "migrate: apply $version"
    $PSQL -1 -f "$file" -c "INSERT INTO schema_migrations (version) VALUES ('$version');"
done

echo "migrate: done"
