CREATE OR REPLACE FUNCTION toggle_favorite(
    mid NUMERIC,
    tid NUMERIC
)
RETURNS NUMERIC AS $$
DECLARE
    row_exists NUMERIC;
BEGIN

    SELECT 1
    INTO row_exists
    FROM favorite
    WHERE member_id = mid and thread_id = tid;

    IF (row_exists > 0) THEN
        DELETE FROM favorite WHERE member_id = mid and thread_id = tid;
        RETURN 0;
    ELSE
        INSERT INTO favorite(member_id, thread_id) VALUES(mid, tid);
        RETURN 1;
    END IF;

END;
$$
LANGUAGE plpgsql