-- +goose Up
-- 00003_add_user_id_to_notes.sql
ALTER TABLE notes
    ADD COLUMN IF NOT EXISTS user_id BIGINT;

ALTER TABLE notes
    ADD CONSTRAINT fk_notes_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
