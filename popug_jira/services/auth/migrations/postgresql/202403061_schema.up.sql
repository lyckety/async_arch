ALTER TABLE users
ADD COLUMN public_id uuid DEFAULT gen_random_uuid() NOT NULL UNIQUE;

ALTER TABLE users
ALTER COLUMN created_at TYPE timestamp without time zone,
ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE users
ALTER COLUMN updated_at TYPE timestamp without time zone;

ALTER TABLE users
ALTER COLUMN deleted_at TYPE timestamp without time zone;
