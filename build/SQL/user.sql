\connect go_api;

-- User Profile -------------------------------------------------------------------------------------------------------------------------------------
-- DROP SEQUENCE IF EXISTS public.seq_usr_profile_id;
CREATE SEQUENCE if not exists public.seq_usr_profile_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_profile;
CREATE TABLE if not exists public.usr_profile (
    id bigint DEFAULT nextval('seq_usr_profile_id' :: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" varchar(100) NOT NULL,
    permissions text [] NOT NULL,
    CONSTRAINT pkey_usr_profile PRIMARY KEY (id),
    CONSTRAINT uni_usr_profile UNIQUE (name)
);

INSERT INTO
    public.usr_profile (id, "name", permissions)
VALUES
    (1, 'ROOT', ARRAY ['*']);

ALTER SEQUENCE IF EXISTS seq_usr_profile_id RESTART WITH 5;

-- User Auth ----------------------------------------------------------------------------------------------------------------------------------------
-- DROP sequence IF EXISTS public.seq_usr_auth_id;
CREATE SEQUENCE if not exists public.seq_usr_auth_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_auth;
CREATE TABLE if not exists public.usr_auth (
    id bigint DEFAULT nextval('seq_usr_auth_id' :: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "status" bool NOT NULL,
    profile_id int8 NOT NULL,
    "token" varchar(255) NULL,
    "password" varchar(255) NULL,
    CONSTRAINT pkey_usr_auth PRIMARY KEY (id),
    CONSTRAINT fk_usr_auth_profile FOREIGN KEY (profile_id) REFERENCES public.usr_profile (id),
    CONSTRAINT uni_usr_auth UNIQUE (token)
);

CREATE INDEX if not exists idx_usr_auth_token ON public.usr_auth USING btree (token);

-- Password: 12345678
INSERT INTO
    public.usr_auth (id, "status", profile_id, token, "password")
VALUES
    (1, true, 1, 'd048aee9-dd65-4ca0-aee7-230c1bf19d8c', '$2a$10$vqkyIvgHRU2sl2FGtlbkNeGFeTsJHQYz18abMJiLlGyJt.Ge99zYy');

ALTER SEQUENCE IF EXISTS seq_usr_auth_id RESTART WITH 5;

-- User ---------------------------------------------------------------------------------------------------------------------------------------------
-- DROP SEQUENCE IF EXISTS public.seq_usr_user_id;
CREATE SEQUENCE if not exists public.seq_usr_user_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_user;
CREATE TABLE if not exists public.usr_user (
    id bigint DEFAULT nextval('seq_usr_user_id' :: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    auth_id int8 NOT NULL,
    "name" varchar(255) NOT NULL,
    mail varchar(50) NOT NULL,
    CONSTRAINT pkey_usr_user PRIMARY KEY (id),
    CONSTRAINT fk_usr_user_auth FOREIGN KEY (auth_id) REFERENCES public.usr_auth (id) ON DELETE CASCADE,
    CONSTRAINT uni_usr_user UNIQUE (mail)
);

CREATE INDEX if not exists idx_users_email ON public.usr_user USING btree (mail);

INSERT INTO
    public.usr_user (id, auth_id, name, mail)
VALUES
    (1, 1, 'Administrator', 'admin@admin.com');

ALTER SEQUENCE IF EXISTS seq_usr_user_id RESTART WITH 5;
