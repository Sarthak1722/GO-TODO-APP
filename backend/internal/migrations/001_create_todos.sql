CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    body TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE
);

---- create above / drop below ----

DROP TABLE todos;