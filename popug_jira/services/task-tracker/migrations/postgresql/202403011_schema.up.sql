--Create user role enums:

CREATE TYPE user_role AS ENUM (
    'administrator',
    'bookkeeper',
    'manager',
    'worker'
);

--Create table users:

CREATE TABLE IF NOT EXISTS users
(
    id uuid PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    role user_role NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
)
WITH (
    OIDS = FALSE
);

--Create task status enums:

CREATE TYPE task_status AS ENUM (
    'opened',
    'completed'
);

--Create table tasks:

CREATE TABLE IF NOT EXISTS tasks
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    description VARCHAR(255) NOT NULL,
    status task_status DEFAULT 'opened',
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
);
