CREATE TABLE chat_group
(
    id    serial PRIMARY KEY,
    topic text
);

CREATE TABLE chat
(
    id         serial PRIMARY KEY,
    member_id  int NOT NULL,
    stamp      timestamp DEFAULT now(),
    chat       text,
    chat_group serial REFERENCES chat_group (id)
);

CREATE TABLE chat_group_member
(
    group_id  serial REFERENCES chat_group (id),
    member_id serial REFERENCES member (id),
    PRIMARY KEY (group_id, member_id)
);

