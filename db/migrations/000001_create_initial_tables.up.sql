CREATE TABLE board_data
(
    id      serial PRIMARY KEY,
    name    text NOT NULL CHECK(name <> ''),  -- name of variable
    value   text NOT NULL CHECK(value <> ''), -- preference value
    UNIQUE(name)
);
INSERT INTO board_data (name,value) VALUES ('total_members',0);
INSERT INTO board_data (name,value) VALUES ('total_threads',0);
INSERT INTO board_data (name,value) VALUES ('total_thread_posts',0);


CREATE TABLE member
(
    id                   serial PRIMARY KEY,
    name                 varchar(25) NOT NULL CHECK(name <> ''),     -- login name
    pass                 char(32) NOT NULL CHECK(pass <> ''),        -- member password md5 hashed
    secret               char(32) NOT NULL CHECK(secret <> ''),      -- secret word for password recovery md5 hashed
    ip                   cidr NOT NULL,                              -- ip of member at last login
    date_joined          timestamp DEFAULT now(),                    -- date of signup
    date_first_post      date,                                       -- the date of the member's first post
    email_signup         text NOT NULL CHECK(email_signup <> ''),    -- email used to sign up
    postalcode           text NOT NULL CHECK(postalcode <> ''),      -- member's postalcode
    total_threads        int DEFAULT 0,                              -- member's total threads created
    total_thread_posts   int DEFAULT 0,                              -- member's total posts
    last_view            timestamp,                                  -- last view of board
    last_post            timestamp,                                  -- last post to board
    last_chat            timestamp,                                  -- last time user chatted
    last_search          timestamp,                                  -- last time user searched
    banned               bool DEFAULT false,                         -- banned user?
    is_admin             bool DEFAULT false,                         -- is admin?
    cookie               char(32),
    session              char(32)
);

CREATE TABLE member_ignore
(
    member_id         int,
    ignore_member_id  int
);

CREATE TABLE member_lurk_unlock
(
    id           serial PRIMARY KEY,
    member_id    int NOT NULL REFERENCES member(id),
    created_at   date NOT NULL DEFAULT now()
);

CREATE TABLE pref_type
(
    id      serial PRIMARY KEY,
    name    text NOT NULL CHECK(name <> ''),
    UNIQUE(name)
);
INSERT INTO pref_type (name) VALUES ('input');
INSERT INTO pref_type (name) VALUES ('checkbox');
INSERT INTO pref_type (name) VALUES ('textarea');

CREATE TABLE pref
(
    id            serial PRIMARY KEY,
    pref_type_id  int NOT NULL REFERENCES pref_type(id),
    name          text NOT NULL CHECK(name <> ''),
    display       text NOT NULL CHECK(display <> ''),
    profile       bool NOT NULL DEFAULT false,
    session       bool NOT NULL DEFAULT false,
    editable      bool NOT NULL DEFAULT true,
    width         int,
    ordering      int,
    UNIQUE(name)
);
INSERT INTO pref VALUES (1,1,'photo','photo url',false,false,true,300,1);
INSERT INTO pref VALUES (2,1,'location','location',true,false,true,200,2);
INSERT INTO pref VALUES (3,1,'email','email',true,false,true,200,3);
INSERT INTO pref VALUES (4,1,'website','website',true,false,true,200,4);
INSERT INTO pref VALUES (5,1,'aim','aim',true,false,true,NULL,5);
INSERT INTO pref VALUES (6,1,'msn','msn',true,false,true,NULL,6);
INSERT INTO pref VALUES (7,1,'yahoo','yahoo',true,false,true,NULL,7);
INSERT INTO pref VALUES (8,1,'gtalk','gtalk',true,false,true,NULL,8);
INSERT INTO pref VALUES (9,1,'jabber','jabber',true,false,true,NULL,9);
INSERT INTO pref VALUES (10,3,'info','info',true,false,true,NULL,10);
INSERT INTO pref VALUES (11,2,'show_email','show email',false,false,true,NULL,12);
INSERT INTO pref VALUES (12,2,'hidemedia','hide media',false,true,true, NULL,13);
INSERT INTO pref VALUES (13,2,'ignore','soft ignore',false,true,true,NULL,14);
INSERT INTO pref VALUES (14,2,'nocollapse','disable collapsing',false,true,true,NULL,19);
INSERT INTO pref VALUES (15,3,'theme','theme',false,true,false,NULL,23);
INSERT INTO pref VALUES (17,2,'nofirstpost','hide firstpost arrow',false,true,true,NULL,15);
INSERT INTO pref VALUES (18,2,'italicread','italicize read posts',false,true,true,NULL,20);
INSERT INTO pref VALUES (19,2,'nopostnumber','hide posts #',false,true,true,NULL,21);
INSERT INTO pref VALUES (20,2,'notabs','hide nav tabs',false,true,true,NULL,22);
INSERT INTO pref VALUES (21,1,'mincollapse','<span class=''small''>min post # to collapse</span>',false,true,true,50,16);
INSERT INTO pref VALUES (22,1,'collapseopen','<span class=''small''># open after collapse (min 1)</span>',false,true,true,50,17);
INSERT INTO pref VALUES (23,1,'externalcss','external css<br/><span class=''small''>(may break color schemes)</span>',false,true,true,300,11);
INSERT INTO pref VALUES (24,1,'uncollapsecount','<span class=''small''># posts to uncollapse</span>',false,true,true,50,18);


CREATE TABLE member_pref
(
    id             serial PRIMARY KEY,
    pref_id        int NOT NULL,
    member_id      int NOT NULL,
    value          text NOT NULL CHECK(value <> '')
);

CREATE TABLE thread
(
    id                 serial PRIMARY KEY,
    member_id          int NOT NULL,                               -- id of member who created thread
    subject            varchar(200) NOT NULL CHECK(subject <> ''), -- subject of thread
    date_posted        timestamp not NULL DEFAULT now(),           -- date thread was created
    first_post_id      int,                                        -- first post id
    posts              int DEFAULT 0,                              -- total posts in a thread
    views              int DEFAULT 0,                              -- total views to thread
    sticky             bool DEFAULT false,                         -- thread sticky flag
    locked             bool DEFAULT false,                         -- thread locked flag
    last_member_id     int NOT NULL,                               -- last member who posted to thread
    date_last_posted   timestamp NOT NULL DEFAULT now(),           -- time last post was entered
    indexed            bool NOT NULL DEFAULT false,                -- has been indexed: for search indexer
    edited             bool NOT NULL DEFAULT false,                -- has been edited: for search indexer
    deleted            bool NOT NULL DEFAULT false,                -- flagged for deletion: for search indexer
    legendary          bool NOT NULL DEFAULT false
);

CREATE TABLE thread_post
(
    id            serial PRIMARY KEY,
    thread_id     int NOT NULL,                  -- thread post belongs to
    date_posted   timestamp DEFAULT now(),       -- time this post was created
    member_id     int NOT NULL,                  -- id of member who created post
    member_ip     cidr NOT NULL,                 -- current ip address of member who created post
    indexed       bool NOT NULL DEFAULT false,   -- has been indexed by search indexer
    edited        bool NOT NULL DEFAULT false,   -- has been edited: for search indexer
    deleted       bool NOT NULL DEFAULT false,   -- flagged for deletion: for search indexer
    body          text                           -- body text of post
);

CREATE TABLE thread_member
(
    member_id	        int NOT NULL,
    thread_id	        int NOT NULL,
    undot                 bool NOT NULL DEFAULT false,
    ignore                bool NOT NULL DEFAULT false,
    date_posted           timestamp,
    last_view_posts       int NOT NULL DEFAULT 0
);


CREATE TABLE message
(
    id                 serial PRIMARY KEY,
    member_id          int NOT NULL,                               -- id of member who created message
    subject            varchar(200) NOT NULL CHECK(subject <> ''), -- subject of message
    first_post_id      int,                                        -- first post id for message
    date_posted        timestamp NOT NULL DEFAULT now(),           -- date message was created
    posts              int DEFAULT 0,                              -- total posts in a message
    views              int DEFAULT 0,                              -- total views to message
    last_member_id     int NOT NULL,                               -- last member to reply
    date_last_posted   timestamp NOT NULL DEFAULT now()
);

CREATE TABLE message_post
(
    id            serial PRIMARY KEY,
    message_id    int NOT NULL,            -- message post belongs to
    date_posted   timestamp DEFAULT now(), -- time message post was created
    member_id     int NOT NULL,            -- id of member who created message post
    member_ip     cidr NOT NULL,           -- current ip address of member who created message post
    body          text                     -- body text of message post
);

CREATE TABLE message_member
(
    member_id	        int NOT NULL,
    message_id	      int NOT NULL,
    date_posted       timestamp,
    last_view_posts   int NOT NULL DEFAULT 0,
    deleted           bool NOT NULL DEFAULT false
);

CREATE TABLE favorite
(
    id          serial PRIMARY KEY,
    member_id   int NOT NULL,       -- member who this watched thread belongs to
    thread_id   int NOT NULL        -- thread member is watching
);