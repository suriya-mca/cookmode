CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT UNIQUE CHECK (char_length(username) BETWEEN 3 AND 30),
    display_name TEXT NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_recipes_updated_at ON recipes;
CREATE TRIGGER trg_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_shopping_lists_updated_at ON shopping_lists;
CREATE TRIGGER trg_shopping_lists_updated_at
    BEFORE UPDATE ON shopping_lists
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

ALTER TABLE recipes DROP CONSTRAINT IF EXISTS recipes_difficulty_check;
ALTER TABLE recipes ADD CONSTRAINT recipes_difficulty_check
    CHECK (difficulty IN ('', 'easy', 'medium', 'hard'));

ALTER TABLE recipes DROP CONSTRAINT IF EXISTS recipes_status_check;
ALTER TABLE recipes ADD CONSTRAINT recipes_status_check
    CHECK (status IN ('draft', 'processing', 'published', 'archived'));

CREATE INDEX IF NOT EXISTS idx_recipes_status_created ON recipes (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_recipes_title_trgm ON recipes USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_recipes_description_trgm ON recipes USING GIN (description gin_trgm_ops);
