set -e

export PGPASSWORD="${POSTGRES_PASSWORD}";

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "postgres" <<-EOSQL
    -- Create the database user if it doesn't exist
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'project_Team29') THEN
            CREATE USER project_Team29 WITH PASSWORD '${POSTGRES_PASSWORD}';
        END IF;
    END
    \$\$;

    -- Create the database if it doesn't exist
    SELECT 'CREATE DATABASE ${POSTGRES_DB}'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${POSTGRES_DB}')\gexec

    -- Grant privileges
    GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO project_Team29;
    ALTER DATABASE ${POSTGRES_DB} OWNER TO project_Team29;
EOSQL