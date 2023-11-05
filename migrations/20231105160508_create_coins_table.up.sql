CREATE TABLE IF NOT EXISTS coins (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL,
year integer NOT NULL,
genres text[] NOT NULL,
price integer NOT NULL
);