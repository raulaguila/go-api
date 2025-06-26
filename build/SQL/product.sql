\connect go_api;

-- DROP SEQUENCE IF EXISTS public.seq_prod_product_id;
CREATE SEQUENCE IF NOT EXISTS public.seq_prod_product_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.prod_product;
CREATE table IF NOT EXISTS public.prod_product (
    id         bigint      DEFAULT nextval('seq_prod_product_id'::regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW()                               NOT NULL,
    updated_at timestamptz DEFAULT NOW()                               NOT NULL,
    "name"     varchar(255)                                            NOT NULL,
    CONSTRAINT pkey_prod_product PRIMARY KEY (id),
    CONSTRAINT uni_prod_product UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_prod_product_name ON public.prod_product USING btree (name);