DROP FUNCTION toggle_ignore(mid NUMERIC, tid NUMERIC);

CREATE OR REPLACE FUNCTION toggle_ignore(
    mid NUMERIC,
    tid NUMERIC
)
RETURNS NUMERIC AS $$
DECLARE
    row_exists NUMERIC;
BEGIN

    SELECT 1
    INTO row_exists
    FROM thread_ignore
    WHERE member_id = mid and ignore_thread_id = tid;

    IF (row_exists > 0) THEN
        DELETE FROM thread_ignore WHERE member_id = mid and ignore_thread_id = tid;
        RETURN 0;
    ELSE
        INSERT INTO thread_ignore(member_id, ignore_thread_id) VALUES(mid, tid);
        RETURN 1;
    END IF;

END;
$$
LANGUAGE plpgsql