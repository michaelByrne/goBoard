SELECT m.id                                                                                      as message,
       m.date_last_posted                                                                        as date_last_posted,
       mem.id,
       mem.name,
       l.id                                                                                      as lastid,
       l.name                                                                                    as lastname,
       m.subject,
       m.posts,
       m.views,
       mp.body
--        (CASE WHEN mm.last_view_posts IS null THEN 0 ELSE mm.last_view_posts END)                 as readbars,
--        (CASE WHEN mm.date_posted IS NOT null AND mm.member_id IS NOT null THEN 't' ELSE 'f' END) as dot
FROM message_member mm
         LEFT JOIN
     message m
     ON
         m.id = mm.message_id
         LEFT JOIN
     member mem
     ON
         mem.id = m.member_id
         LEFT JOIN
     member l
     ON
         l.id = m.last_member_id
         LEFT JOIN
     message_post mp
     ON
         mp.id = m.first_post_id
WHERE mm.member_id = $1
  AND mm.deleted IS false