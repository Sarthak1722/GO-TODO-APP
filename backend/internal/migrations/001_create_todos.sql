CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    body TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos (user_id, id);

---- create above / drop below ----

DROP TABLE todos;
