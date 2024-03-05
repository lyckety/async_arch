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
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    role user_role NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
)
WITH (
    OIDS = FALSE
);
