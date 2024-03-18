--Create user role enums:

CREATE TYPE user_role AS ENUM (
    'unspecified',
    'administrator',
    'bookkeeper',
    'manager',
    'worker'
);

--Create table users:

CREATE TABLE IF NOT EXISTS users
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    public_id uuid NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    username VARCHAR(50),
    email VARCHAR(100),
    role user_role DEFAULT 'unspecified',
    balance integer DEFAULT 0,
    created_at timestamp without time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
)
WITH (
    OIDS = FALSE
);

--Create table tasks:

CREATE TABLE IF NOT EXISTS tasks
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    public_id uuid NOT NULL UNIQUE,
    description VARCHAR(255),
    cost_assign smallint NOT NULL,
    cost_complete smallint NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
)
WITH (
    OIDS = FALSE
);

--Create billing_cycle_status enums:

CREATE TYPE billing_cycle_status AS ENUM (
    'opened',
    'started',
    'completed'
);

--Create table billing_cycles:

CREATE TABLE IF NOT EXISTS billing_cycles
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_public_id uuid NOT NULL,
    status billing_cycle_status NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
)
WITH (
    OIDS = FALSE
);

--Create transaction_type enums:

CREATE TYPE transaction_type AS ENUM (
    'task_assigned',
    'task_completed',
    'payment'
);

--Create table transactions:

CREATE TABLE IF NOT EXISTS transactions
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    public_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_public_id uuid NOT NULL,
    task_public_id uuid NOT NULL,
    type transaction_type NOT NULL,
    debit integer DEFAULT 0 NOT NULL,
    credit integer DEFAULT 0 NOT NULL,
    date_time timestamp without time zone,
    created_at timestamp without time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
)
WITH (
    OIDS = FALSE
);
