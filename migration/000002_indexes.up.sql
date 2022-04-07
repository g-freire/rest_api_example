-- CLASS

CREATE INDEX IF NOT EXISTS  "idx_class_by_dates"
    ON public.class USING btree
        (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
    INCLUDE(creation_time, name, capacity)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS "idx_class_by_name"
    ON public.class USING brin(name)
    TABLESPACE pg_default;

-- covering index -- too much space in big scale
-- CREATE INDEX "idx_class_by_id"
--     ON public.class USING btree(id)
--     INCLUDE(creation_time, name, start_date, end_date, capacity)
--     TABLESPACE pg_default;

-- MEMBER
CREATE INDEX IF NOT EXISTS "idx_member_by_name"
    ON public.member USING brin(name)
    TABLESPACE pg_default;


-- -- BOOKING
CREATE INDEX IF NOT EXISTS "idx_booking_by_fks"
ON public.booking USING btree(class_id, member_id)
TABLESPACE pg_default;


