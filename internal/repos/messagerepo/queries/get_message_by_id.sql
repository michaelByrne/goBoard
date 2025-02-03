SELECT
    m.id,
    m.date_last_posted,
    mem.id,
    mem.name,
    m.last_member_id,
    l.name,
    m.subject,
    m.posts,
    m.views,
    mp.body
FROM message_member mm
LEFT JOIN message m ON m.id = mm.message_id
LEFT JOIN member mem ON m.member_id = mem.id
LEFT JOIN member l ON l.id = m.last_member_id
LEFT JOIN message_post mp ON mp.id = m.first_post_id
WHERE mm.member_id = $1
AND m.id = $2
AND mm.deleted IS FALSE