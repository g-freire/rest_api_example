CREATE TABLE IF NOT EXISTS class(
    id serial PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(80) NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    capacity SMALLINT NULL
);


CREATE INDEX "idx_class_by_id"
    ON public.class USING btree(id)
    INCLUDE(created_at, name, start_date, end_date, capacity)
    TABLESPACE pg_default;


CREATE INDEX "idx_class_by_dates"
    ON public.class USING btree
        (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
    INCLUDE(created_at, name, capacity)
    TABLESPACE pg_default;