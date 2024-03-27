CREATE OR REPLACE FUNCTION toggle_ignore(
    member_id NUMERIC, 
    thread_id NUMERIC
) 
RETURNS NUMERIC AS $$
DECLARE
    row_exists NUMERIC;
BEGIN

    SELECT 1 
    INTO row_exists 
    FROM thread_ignore 
    WHERE member_id = member_id and ignore_thread_id = thread_id;

    IF (row_exists > 0) THEN
        DELETE FROM thread_ignore WHERE member_id = member_id and ignore_thread_id = thread_id;
        RETURN 0;
    ELSE
        INSERT INTO thread_ignore(member_id, ignore_thread_id) VALUES(member_id, thread_id);
        RETURN 1;
    END IF;

END; 
$$
LANGUAGE plpgsql