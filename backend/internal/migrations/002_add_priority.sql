ALTER TABLE todos ADD COLUMN priority INTEGER DEFAULT 1;

---- create above / drop below ----

ALTER TABLE todos DROP COLUMN priority;