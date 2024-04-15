with collapsible AS (SELECT tp.id,
                            tp.date_posted,
                            m.id                                                         AS member_id,
                            m.name,
                            tp.body,
                            tp.member_ip,
                            t.subject,
                            t.id                                                         AS thread_id,
                            m.is_admin,
                            ROW_NUMBER() OVER (ORDER BY tp.date_posted DESC, tp.id DESC) AS row_num
                     FROM thread_post tp
                              LEFT JOIN member m on m.id = tp.member_id
                              LEFT JOIN thread t on t.id = tp.thread_id
                     WHERE tp.thread_id = 60
                       AND m.id NOT IN (SELECT m.id
                                        FROM member_ignore mi
                                                 LEFT JOIN
                                             member m ON m.id = mi.ignore_member_id
                                        WHERE mi.member_id = 2
                                        ORDER BY m.name)),
     count_less_than_start AS (SELECT COUNT(*) AS count_less
                               FROM collapsible
                               WHERE row_num < 5)
SELECT id,
       date_posted,
       member_id,
       name,
       body,
       member_ip,
       subject,
       thread_id,
       is_admin,
       row_num,
       (SELECT count_less FROM count_less_than_start) AS count_less_than_start
FROM collapsible as cp
WHERE row_num BETWEEN 5 AND 10
ORDER BY row_num;