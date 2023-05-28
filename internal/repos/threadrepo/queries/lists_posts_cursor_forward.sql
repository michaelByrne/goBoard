SELECT tp.id,
       tp.date_posted,
       m.id as member_id,
       m.name,
       tp.body,
       tp.member_ip,
       t.subject,
       t.id as thread_id,
       m.is_admin
FROM thread_post tp
         LEFT JOIN
     member m
     ON
         m.id = tp.member_id
         LEFT JOIN
     thread t
     ON
         t.id = tp.thread_id
WHERE tp.thread_id = $2
AND tp.date_posted > $3
ORDER BY tp.date_posted ASC, tp.id ASC
LIMIT $1
