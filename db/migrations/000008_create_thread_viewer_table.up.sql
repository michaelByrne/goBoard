DROP TABLE IF EXISTS thread_ignore;

CREATE TABLE thread_viewer (
    member_id int,
    thread_id int,
    ignored bool DEFAULT false,
    dotted bool DEFAULT false
);

CREATE UNIQUE INDEX tv_index ON thread_viewer(member_id, thread_id);
ALTER TABLE thread_viewer ADD FOREIGN KEY (member_id) REFERENCES member(id);
ALTER TABLE thread_viewer ADD FOREIGN KEY (thread_id) REFERENCES thread(id);