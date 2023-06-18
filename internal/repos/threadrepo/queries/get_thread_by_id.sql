SELECT id,
       subject,
       t.date_posted,
       t.member_id,
       views,
       (CASE
            WHEN tm.date_posted IS NOT null AND tm.undot IS false AND tm.member_id IS NOT null THEN true
            ELSE false END) as dot,
       tm.ignore
FROM thread t
         LEFT OUTER JOIN thread_member tm on t.id = tm.thread_id AND tm.member_id = $2
WHERE id = $1