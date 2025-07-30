CREATE TABLE posts (
    id TEXT PRIMARY KEY DEFAULT nanoid('post_', 25);
    title VARCHAR(255) NOT NULL;
    content TEXT NOT NULL;
    is_cool BOOLEAN;
    created_at TIMESTAMP DEFAULT NOW();
);
