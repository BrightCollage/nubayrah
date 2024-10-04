CREATE TABLE library (
    id integer PRIMARY KEY,
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
    uid test
);