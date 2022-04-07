-- DDF
CREATE TABLE IF NOT EXISTS class(
    id serial PRIMARY KEY,
    name VARCHAR(80) NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    capacity SMALLINT NULL,  --add booked
    creation_time TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS member(
    id serial PRIMARY KEY,
    name VARCHAR(80) NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT NOW()
);

-- associative entity
CREATE TABLE IF NOT EXISTS booking(
    id serial PRIMARY KEY, -- surrogate key to allow many bookings to be made for a single class
    class_id    int NOT NULL REFERENCES class (id), --FK
    member_id   int NOT NULL REFERENCES member(id), --FK
    date TIMESTAMP NOT NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT NOW(),
    --UNIQUE(class_id, member_id, date), -- a member cannot book the same class twice on the same date
    UNIQUE(member_id, date) -- a member cannot book a class twice on the same date
);

-- INDEXES
CREATE INDEX "idx_class_by_id"
    ON public.class USING btree(id)
    INCLUDE(creation_time, name, start_date, end_date, capacity)
    TABLESPACE pg_default;


CREATE INDEX "idx_class_by_dates"
    ON public.class USING btree
        (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
    INCLUDE(creation_time, name, capacity)
    TABLESPACE pg_default;