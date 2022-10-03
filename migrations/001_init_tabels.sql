-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users
(
    user_id  SERIAL PRIMARY KEY,
    username varchar(255) NOT NULL UNIQUE,
    fullname varchar(255) NOT NULL,
    password varchar(75)  NOT NULL,
    email    varchar(255) NOT NULL,
    role     smallint     NOT NULL
);

CREATE TABLE IF NOT EXISTS vote
(
    vote_id     SERIAL PRIMARY KEY,
    title       varchar(256) NOT NULL UNIQUE,
    slug        varchar(256) NOT NULL UNIQUE,
    vote_type   smallint     NOT NULL,
    vote_status smallint     NOT NULL
);

CREATE TABLE IF NOT EXISTS vote_options
(
    opts_id     SERIAL PRIMARY KEY,
    user_id     integer,
    option_name varchar(256) NOT NULL,
    all_votes   int          NOT NULL,
    quiz_answer boolean      NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS user_votes
(
    id      SERIAL PRIMARY KEY,
    vote_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT fk_vote_id FOREIGN KEY (vote_id)
        REFERENCES vote (vote_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users (user_id)
);

---- create above / drop below ----
DROP TABLE user_votes;
DROP TABLE vote_options;
DROP TABLE vote;
DROP TABLE users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
