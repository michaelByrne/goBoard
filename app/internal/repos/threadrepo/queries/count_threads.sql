SELECT
    b.value::int as total_threads
FROM
    board_data b
WHERE
    b.name = 'total_threads'
