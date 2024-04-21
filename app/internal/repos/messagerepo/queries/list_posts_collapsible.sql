WITH collapsed_posts AS (
    SELECT
        mp.id,
        mp.date_posted,
        m.id as member_id,
        m.name,
        mp.body,
        mp.member_ip,
        mg.subject,
        mg.id as message_id,
        ROW_NUMBER() OVER (ORDER BY mp.date_posted DESC) AS original_row_num,
        COUNT(*) OVER() - $3 AS total_collapsed
    FROM message_post mp
    LEFT JOIN member m ON mp.member_id = m.id
    LEFT JOIN message mg ON mp.message_id = mg.id
    WHERE mp.message_id = $1
    AND EXISTS (
        SELECT mm.member_id
        FROM message_member mm
        WHERE mm.message_id = mp.message_id
          AND mm.member_id = $2
    )
)
SELECT
    id,
    date_posted,
    member_id,
    name,
    body,
    member_ip,
    subject,
    message_id,
    $3 - (cp.original_row_num - (cp.total_collapsed - $3)) + $3 AS row_num,
    total_collapsed
FROM collapsed_posts cp
WHERE cp.original_row_num <= $3
ORDER BY cp.date_posted ASC