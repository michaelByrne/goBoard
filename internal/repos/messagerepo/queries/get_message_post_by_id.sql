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
WHERE mp.id = $1