-- name: InsertChat :one
INSERT INTO chat (member_id, chat, chat_group)
VALUES ($1, $2, $3)
RETURNING *;

-- name: InsertChatGroup :one
INSERT INTO chat_group (topic)
VALUES ($1)
RETURNING *;

-- name: InsertChatGroupMember :one
INSERT INTO chat_group_member (group_id, member_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetChatGroupsForMember :many
SELECT cg.id, cg.topic
FROM chat_group cg
JOIN chat_group_member cgm ON cg.id = cgm.group_id
WHERE member_id = $1;

-- name: GetChatsForGroup :many
SELECT c.id, c.member_id, c.stamp, c.chat
FROM chat c
WHERE c.chat_group = $1;