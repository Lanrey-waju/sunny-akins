\set pwd `echo $SUNNY_PASSWORD`
\set mig_pwd `echo $MIGRATIONS_PASSWORD`

\echo :pwd
\echo :mig_pwd

SET app.migrations_password = :mig_pwd;
SET app.app_password = :pwd;

DO $$
BEGIN
    -- Create migrations user if it doesn't exist
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'migrations') THEN
        EXECUTE format('CREATE USER migrations WITH PASSWORD %L', current_setting('app.migrations_password'));
        ALTER USER migrations WITH SUPERUSER;
    END IF;
    -- Create sunny user if it doesn't exi`t
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'sunny') THEN
        EXECUTE format('CREATE USER sunny WITH PASSWORD %L', current_setting('app.app_password'));
        ALTER USER migrations WITH SUPERUSER;
    END IF;
    -- Grant privileges for sunny
END
$$;

\c sunny
-- Create schema if not exists
CREATE SCHEMA IF NOT EXISTS sunny_schema;

-- Grant connect privilege
GRANT CONNECT ON DATABASE sunny TO sunny;

-- Grant privileges on both public and sunny_schema
GRANT USAGE ON SCHEMA public TO sunny;
GRANT USAGE ON SCHEMA sunny_schema TO sunny;

-- Set default privileges for future objects created by migrations user
ALTER DEFAULT PRIVILEGES FOR ROLE migrations IN SCHEMA public 
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO sunny;
ALTER DEFAULT PRIVILEGES FOR ROLE migrations IN SCHEMA public 
    GRANT USAGE, SELECT ON SEQUENCES TO sunny;

ALTER DEFAULT PRIVILEGES FOR ROLE migrations IN SCHEMA sunny_schema 
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO sunny;
ALTER DEFAULT PRIVILEGES FOR ROLE migrations IN SCHEMA sunny_schema 
    GRANT USAGE, SELECT ON SEQUENCES TO sunny;

