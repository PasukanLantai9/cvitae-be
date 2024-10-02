CREATE TABLE resumes (
    id VARCHAR(255) primary key,
    name VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL
)