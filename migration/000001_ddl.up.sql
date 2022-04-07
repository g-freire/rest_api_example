-- DDF
CREATE TABLE IF NOT EXISTS class(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(80) NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    capacity SMALLINT NULL,  --add booked
    creation_time TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS member(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,    name VARCHAR(80) NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT NOW()
);

-- associative entity
CREATE TABLE IF NOT EXISTS booking(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- surrogate key to allow many bookings to be made for a single class
    class_id    int NOT NULL REFERENCES class (id), --FK
    member_id   int NOT NULL REFERENCES member(id), --FK
    date TIMESTAMP NOT NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT NOW(),
    --UNIQUE(class_id, member_id, date), -- a member cannot book the same class twice on the same date
    UNIQUE(member_id, date) -- a member cannot book a class twice on the same date
);

