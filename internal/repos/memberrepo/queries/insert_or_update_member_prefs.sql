INSERT INTO member_pref (pref_id, member_id, value) VALUES ((SELECT id FROM pref WHERE name = $1), $2, $3) ON CONFLICT (id)
    DO UPDATE SET value = $3