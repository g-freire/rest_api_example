-- CLASS
CREATE INDEX "idx_class_by_id"
    ON public.class USING btree(id)
    INCLUDE(creation_time, name, start_date, end_date, capacity)
    TABLESPACE pg_default;


CREATE INDEX "idx_class_by_dates"
    ON public.class USING btree
        (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
    INCLUDE(creation_time, name, capacity)
    TABLESPACE pg_default;

-- MEMBER
-- CREATE INDEX "idx_class_by_id"
--     ON public.class USING btree(id)
--     INCLUDE(creation_time, name, start_date, end_date, capacity)
--     TABLESPACE pg_default;
--
--
-- CREATE INDEX "idx_class_by_dates"
--     ON public.class USING btree
--         (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
--     INCLUDE(creation_time, name, capacity)
--     TABLESPACE pg_default;
--
-- -- BOOKING
-- CREATE INDEX "idx_class_by_id"
--     ON public.class USING btree(id)
--     INCLUDE(creation_time, name, start_date, end_date, capacity)
--     TABLESPACE pg_default;
--
--
-- CREATE INDEX "idx_class_by_dates"
--     ON public.class USING btree
--         (start_date DESC NULLS LAST, end_date DESC NULLS LAST)
--     INCLUDE(creation_time, name, capacity)
--     TABLESPACE pg_default;