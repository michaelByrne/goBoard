WITH collapsed_posts AS (
    SELECT
        tp.id,
        tp.date_posted,
        m.id AS member_id,
        m.name,
        tp.body,
        tp.member_ip,
        t.subject,
        t.id AS thread_id,
        m.is_admin,
        ROW_NUMBER() OVER (ORDER BY tp.date_posted DESC, tp.id DESC) AS original_row_num, -- Calculate original row numbers
        COUNT(*) OVER() - $3 AS total_collapsed -- Calculate the total number of collapsed posts
    FROM
        thread_post tp
    LEFT JOIN
        member m ON m.id = tp.member_id
    LEFT JOIN
        thread t ON t.id = tp.thread_id
    WHERE
        tp.thread_id = COALESCE($1, tp.thread_id)
    AND m.id NOT IN (
        SELECT
            m.id
        FROM
            member_ignore mi
        LEFT JOIN
            member m ON m.id = mi.ignore_member_id
        WHERE
            mi.member_id = $2
        ORDER BY
            m.name
    )
)
SELECT
    cp.id,
    cp.date_posted,
    cp.member_id,
    cp.name,
    cp.body,
    cp.member_ip,
    cp.subject,
    cp.thread_id,
    cp.is_admin,
    $3 - (cp.original_row_num - (cp.total_collapsed - $3)) + $3 AS row_num, -- Adjust row numbers based on visible posts
    cp.total_collapsed -- Return the total number of displayed posts
FROM
    collapsed_posts cp
WHERE
    cp.original_row_num <= $3 -- Filter visible posts
ORDER BY
    cp.date_posted ASC, cp.id ASC; -- Order by date_posted and id to ensure consistent ordering of posts with the same date_posted
