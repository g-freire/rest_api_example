-- SEED - SERIES GENERATOR

-- eg. creates 10M rows on each table

-- CONST
set session const.rows = 10000000;
set session const.booking.id.range = 2;

-- CLASS
INSERT INTO class(name, start_date, end_date, capacity)
SELECT md5(RANDOM()::VARCHAR(80)),
       (NOW() + (random() * (NOW()+'1000 days' - NOW())) + '30 days'),
       (NOW() + (random() * (NOW()+'1000 days' - NOW())) + '30 days'),
       (RANDOM()::SMALLINT)
FROM GENERATE_SERIES(1, current_setting('const.rows')::int);

-- MEMBER
INSERT INTO member(name)
SELECT md5(RANDOM()::VARCHAR(80))
FROM GENERATE_SERIES(1, current_setting('const.rows')::int);

-- BOOKING
INSERT INTO booking(class_id, member_id, date)
SELECT (RANDOM() * current_setting('const.booking.id.range')::int + 1::BIGINT),
       (RANDOM() * current_setting('const.booking.id.range')::int + 1::BIGINT),
       (NOW() + (random() * (NOW()+'1000 days' - NOW())) + '30 days')
FROM GENERATE_SERIES(1, current_setting('const.rows')::int);
