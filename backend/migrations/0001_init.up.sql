-- 0001_init: baseline. Feature tables go in later migrations (make new-migration name=...).

-- gen_random_uuid() is built into Postgres 13+; pgcrypto kept for older versions.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Reusable trigger function: keeps updated_at current on every UPDATE.
-- Usage on a new table:
--   CREATE TRIGGER trg_<table>_updated_at BEFORE UPDATE ON <table>
--   FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
