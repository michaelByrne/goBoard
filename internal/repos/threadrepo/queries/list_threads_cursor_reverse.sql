SELECT *
FROM (SELECT t.id                 as thread,
             t.date_last_posted,
             t.date_posted,
             m.id,
             m.name,
             l.id                 as lastid,
             l.name               as lastname,
             t.subject,
             t.posts,
             t.views,
             tp.body,
             t.sticky,
             t.locked,
             t.legendary,
             (CASE
                  WHEN tm.date_posted IS NOT null AND tm.undot IS false AND tm.member_id IS NOT null THEN true
                  ELSE false END) as dot
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
               LEFT OUTER JOIN thread_member tm
                               ON tm.thread_id = t.id AND tm.member_id = $3
      WHERE t.sticky IS false
        AND t.date_last_posted >= $1
        AND m.id NOT IN (SELECT m.id
                         FROM member_ignore mi
                                  LEFT JOIN
                              member m
                              ON
                                  m.id = mi.ignore_member_id
                         WHERE mi.member_id = $3
                         ORDER BY m.name)
      AND tm.ignore IS NOT true
      ORDER BY t.date_last_posted ASC
      LIMIT $2 + 1) AS pagination
ORDER BY date_last_posted DESC;
