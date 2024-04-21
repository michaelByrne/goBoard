CREATE TABLE message_viewer (
    member_id INT,
    message_id INT,
    dotted bool DEFAULT true,
    last_viewed TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE UNIQUE INDEX mv_index ON message_viewer (member_id, message_id);
ALTER TABLE message_viewer ADD FOREIGN KEY (member_id) REFERENCES member(id);
ALTER TABLE message_viewer ADD FOREIGN KEY (message_id) REFERENCES message(id);
