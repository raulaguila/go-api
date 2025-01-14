\connect go_api;

-- DROP sequence IF EXISTS public.users_auth_id_seq;

CREATE SEQUENCE if not exists public.users_auth_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1
    NO CYCLE;

-- DROP SEQUENCE IF EXISTS public.users_id_seq;

CREATE SEQUENCE if not exists public.users_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1
    NO CYCLE;

-- DROP SEQUENCE IF EXISTS public.users_profile_id_seq;

CREATE SEQUENCE if not exists public.users_profile_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1
    NO CYCLE;

-- DROP TABLE public.users_profile;

CREATE TABLE if not exists public.users_profile
(
    id          bigint      DEFAULT nextval('users_profile_id_seq'::regclass) NOT NULL,
    created_at  timestamptz default NOW()                                     NOT NULL,
    updated_at  timestamptz default NOW()                                     NOT NULL,
    "name"      varchar(100)                                                  NOT NULL,
    permissions text[]                                                        NOT NULL,
    CONSTRAINT uni_users_profile_name UNIQUE (name),
    CONSTRAINT users_profile_pkey PRIMARY KEY (id)
);

INSERT INTO public.users_profile (name, permissions)
VALUES ('ROOT', ARRAY ['*']);


-- DROP TABLE public.users;

CREATE TABLE if not exists public.users
(
    id         bigint      DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
    created_at timestamptz default NOW()                             NOT NULL,
    updated_at timestamptz default NOW()                             NOT NULL,
    "name"     varchar(90)                                           NOT NULL,
    mail       varchar(50)                                           NOT NULL,
    CONSTRAINT uni_users_mail UNIQUE (mail),
    CONSTRAINT users_pkey PRIMARY KEY (id)
--     CONSTRAINT fk_users_auth FOREIGN KEY (auth_id) REFERENCES public.users_auth (id) ON DELETE CASCADE
);

-- CREATE INDEX if not exists idx_users_auth_id ON public.users USING btree (auth_id);
CREATE INDEX if not exists idx_users_email ON public.users USING btree (mail);

INSERT INTO public.users (name, mail)
VALUES ('Administrator', 'admin@admin.com');


-- DROP TABLE public.users_auth;

CREATE TABLE if not exists public.users_auth
(
    id         bigint      DEFAULT nextval('users_auth_id_seq'::regclass) NOT NULL,
    created_at timestamptz default NOW()                                  NOT NULL,
    updated_at timestamptz default NOW()                                  NOT NULL,
    user_id    int8                                                       NOT NULL,
    status     bool                                                       NOT NULL,
    profile_id int8                                                       NOT NULL,
    "token"    varchar(255)                                               NULL,
    "password" varchar(255)                                               NULL,
    CONSTRAINT uni_users_auth_token UNIQUE (token),
    CONSTRAINT users_auth_pkey PRIMARY KEY (id),
    CONSTRAINT fk_users_auth_profile FOREIGN KEY (profile_id) REFERENCES public.users_profile (id),
    CONSTRAINT fk_users_auth_user FOREIGN KEY (user_id) REFERENCES public.users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX if not exists idx_users_auth_profile_id ON public.users_auth USING btree (profile_id);
CREATE INDEX if not exists idx_users_auth_token ON public.users_auth USING btree (token);

-- Password: 12345678
INSERT INTO public.users_auth (user_id, status, profile_id, token, "password")
VALUES (1, true, 1, 'd048aee9-dd65-4ca0-aee7-230c1bf19d8c',
        '$2a$10$vqkyIvgHRU2sl2FGtlbkNeGFeTsJHQYz18abMJiLlGyJt.Ge99zYy');
