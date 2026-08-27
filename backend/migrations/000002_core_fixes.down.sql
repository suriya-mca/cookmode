DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TRIGGER IF EXISTS trg_shopping_lists_updated_at ON shopping_lists;
DROP TRIGGER IF EXISTS trg_recipes_updated_at ON recipes;
DROP FUNCTION IF EXISTS set_updated_at();

DROP INDEX IF EXISTS idx_recipes_description_trgm;
DROP INDEX IF EXISTS idx_recipes_title_trgm;
DROP INDEX IF EXISTS idx_recipes_status_created;

ALTER TABLE recipes DROP CONSTRAINT IF EXISTS recipes_status_check;
ALTER TABLE recipes ADD CONSTRAINT recipes_status_check
    CHECK (status IN ('draft', 'processing', 'published', 'archived'));

ALTER TABLE recipes DROP CONSTRAINT IF EXISTS recipes_difficulty_check;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pg_trgm";
