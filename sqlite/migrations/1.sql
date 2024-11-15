CREATE TABLE library (
    id text PRIMARY KEY,
    filepath text NOT NULL,
    title text NOT NULL,
    titleSort text,
    author text NOT NULL,
    authorSort text,
    language text,
    series text,
    seriesNum text,
    subjects text,
    isbn text,
    publisher text,
    pubDate text,
    rights text,
    contributors text,
    description text,
    uid text
);

CREATE TABLE dbmetadata (
    id int PRIMARY KEY,
    userVersion int
);

INSERT INTO dbmetadata (id, userVersion) VALUES (0, 0);