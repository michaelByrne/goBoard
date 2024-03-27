SELECT t.id,
       t.subject,
       t.date_posted,
       t.member_id,
       t.views,
       COALESCE(tv.dotted, false) as dot,
       COALESCE(tv.undot, false) as undot,
       COALESCE(tv.ignored, false) as ignore,
       CASE WHEN f.thread_id IS NOT NULL THEN true ELSE false END as favorite
FROM thread t
LEFT JOIN thread_viewer tv ON t.id = tv.thread_id AND tv.member_id = $2
LEFT OUTER JOIN favorite f ON t.id = f.thread_id AND f.member_id = $2
WHERE t.id = $1;