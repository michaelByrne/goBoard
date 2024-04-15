SELECT
    tp.id,
    tp.date_posted,
    m.id as member_id,
    m.name,
    tp.body,
    tp.member_ip,
    t.subject,
    t.id as thread_id,
    m.is_admin
FROM
    thread_post tp
LEFT JOIN
    member m
ON
    m.id=tp.member_id
LEFT JOIN
    thread t
ON
    t.id = tp.thread_id
WHERE
    tp.thread_id = coalesce($3, tp.thread_id)
AND m.id NOT IN (
    SELECT
        m.id
    FROM
        member_ignore mi
            LEFT JOIN
        member m
        ON
                m.id = mi.ignore_member_id
    WHERE
            mi.member_id=$4
    ORDER BY
        m.name
    )
ORDER BY tp.date_posted ASC, tp.id ASC
OFFSET $2
LIMIT $1