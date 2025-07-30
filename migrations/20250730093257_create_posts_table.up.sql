CREATE TABLE posts (
    id nanoid%!(EXTRA string=post) PRIMARY KEY;
    title string NOT NULL;
    content text NOT NULL;
    is_cool bool;
    created_at datetime DEFAULT NOW();
);
