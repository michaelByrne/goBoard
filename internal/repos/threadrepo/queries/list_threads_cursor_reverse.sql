SELECT *
FROM (SELECT t.id   as thread,
             t.date_last_posted,
             t.date_posted,
             m.id,
             m.name,
             l.id   as lastid,
             l.name as lastname,
             t.subject,
             t.posts,
             t.views,
             tp.body,
             t.sticky,
             t.locked,
             t.legendary
      FROM thread t
               LEFT JOIN
           member m
           ON
               m.id = t.member_id
               LEFT JOIN
           member l
           ON
               l.id = t.last_member_id
               LEFT JOIN
           thread_post tp
           ON
               tp.id = t.first_post_id
      WHERE t.sticky IS false
        AND t.date_last_posted >= $1
      ORDER BY t.date_last_posted ASC
      LIMIT $2+1) AS pagination
ORDER BY date_last_posted DESC ;
