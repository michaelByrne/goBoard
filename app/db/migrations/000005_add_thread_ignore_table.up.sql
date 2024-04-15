CREATE TABLE thread_ignore
(
    member_id         int,
    ignore_thread_id  int
);

--start member_ignore
CREATE UNIQUE INDEX mi_ti_index ON thread_ignore(member_id, ignore_thread_id);
ALTER TABLE thread_ignore ADD FOREIGN KEY (member_id) REFERENCES member(id);
ALTER TABLE thread_ignore ADD FOREIGN KEY (ignore_thread_id) REFERENCES thread(id);