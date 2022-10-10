-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users
(
    user_id  SERIAL PRIMARY KEY,
    username varchar(255) NOT NULL UNIQUE,
    fullname varchar(255) NOT NULL,
    password varchar(75)  NOT NULL,
    email    varchar(255) NOT NULL UNIQUE,
    role     smallint     NOT NULL DEFAULT 1
);

-- every vote has a creator user
CREATE TABLE IF NOT EXISTS votes
(
    vote_id     SERIAL PRIMARY KEY,
    user_id     integer,
    title       varchar(256) NOT NULL UNIQUE,
    slug        varchar(256) NOT NULL UNIQUE,
    vote_type   smallint     NOT NULL DEFAULT 1,
    vote_status smallint     NOT NULL DEFAULT 2,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vote_options
(
    opts_id     SERIAL PRIMARY KEY,
    vote_id     integer,
    option_name varchar(256) NOT NULL,
    all_votes   int          NOT NULL,
    quiz_answer boolean      NOT NULL,

    CONSTRAINT fk_vote_id FOREIGN KEY (vote_id)
        REFERENCES votes (vote_id)
);

CREATE TABLE IF NOT EXISTS user_votes
(
    id            SERIAL PRIMARY KEY,
    vote_id       integer      NOT NULL,
    user_id       integer      NOT NULL,
    user_username varchar(256) NOT NULL,
    opts_id       integer      NOT NULL,
    CONSTRAINT fk_vote_id FOREIGN KEY (vote_id)
        REFERENCES votes (vote_id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users (user_id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT fk_opts_id FOREIGN KEY (opts_id)
        REFERENCES vote_options (opts_id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT fk_user_username FOREIGN KEY (user_username)
        REFERENCES users (username) ON DELETE CASCADE DEFERRABLE
);


-- SET CONSTRAINTS ALL DEFERRED;
---- create above / drop below ----
DROP TABLE IF EXISTS user_votes;
DROP TABLE IF EXISTS vote_options;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
