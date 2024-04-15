SELECT mp.id,
       mp.date_posted,
       mem.id as member_id,
       mem.name,
       mp.body,
       mp.member_ip,
       m.subject,
       m.id   as message_id
FROM message_post mp
         LEFT JOIN
     member mem
     ON
         mem.id = mp.member_id
         LEFT JOIN
     message m
     ON
         m.id = mp.message_id
WHERE mp.message_id = $1
  AND EXISTS (SELECT mm.member_id
              FROM message_member mm
              WHERE mm.message_id = mp.message_id
                AND mm.member_id = $2)
LIMIT $3