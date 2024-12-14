\connect go_api;

-- DROP SEQUENCE IF EXISTS public.product_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.product_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

CREATE table IF NOT EXISTS public.product (
    id bigint DEFAULT nextval('product_id_seq'::regclass) NOT NULL,
    created_at timestamptz NOT NULL default NOW(),
    updated_at timestamptz NOT NULL default NOW(),
    "name" varchar(100) NOT NULL,
    CONSTRAINT product_pkey PRIMARY KEY (id),
    CONSTRAINT uni_product_name UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_product_name ON public.product USING btree (name);