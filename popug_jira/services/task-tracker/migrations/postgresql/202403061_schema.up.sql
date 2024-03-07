ALTER TABLE tasks
ADD COLUMN public_id uuid DEFAULT gen_random_uuid() NOT NULL UNIQUE;

UPDATE tasks
SET public_id = id;

ALTER TABLE tasks
ALTER COLUMN created_at TYPE timestamp without time zone,
ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE tasks
ALTER COLUMN updated_at TYPE timestamp without time zone;

ALTER TABLE tasks
ALTER COLUMN deleted_at TYPE timestamp without time zone;
