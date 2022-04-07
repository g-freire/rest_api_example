-- LIST INDEXES
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY
    tablename,
    indexname;

-- DETECT MISSING INDEXES
-- finds large tables that have been used frequently in a sequential scan
SELECT
   schemaname,
   relname,
   seq_scan,
   seq_tup_read,
   seq_tup_read / seq_scan AS avg, idx_scan
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 10;

-- FINDS USELESS INDEXES
--how often an index was used and how much space has been wasted for each index
SELECT
   schemaname,
   relname,
   indexrelname,
   idx_scan,
   pg_size_pretty(pg_relation_size(indexrelid)) AS idx_size
FROM pg_stat_user_indexes;

-- FIND TOP TIME-CONSUMING QUERIES
SELECT
   round((100 * total_time / sum(total_time) OVER ())::numeric, 2) percent,
   round(total_time::numeric, 2) AS total,
   calls,
   round(mean_time::numeric, 2) AS mean,
   substring(query, 1, 40)
FROM pg_stat_statements
ORDER BY total_time
DESC LIMIT 10;

