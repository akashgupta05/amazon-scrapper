CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE products (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    link text,
    product text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
ALTER TABLE ONLY products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);