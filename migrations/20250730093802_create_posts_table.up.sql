CREATE TABLE posts (
    id SELECT NANOID('post', 25) PRIMARY KEY;
    title VARCHAR(255) NOT NULL;
    content TEXT NOT NULL;
    is_cool BOOLEAN;
    created_at TIMESTAMP DEFAULT NOW();
);
