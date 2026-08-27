CREATE TABLE follows (
    follower_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_follows_following_id ON follows (following_id);
CREATE INDEX idx_follows_follower_id ON follows (follower_id);

CREATE TABLE saves (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, recipe_id)
);

CREATE INDEX idx_saves_recipe_id ON saves (recipe_id);
CREATE INDEX idx_saves_user_id ON saves (user_id);

CREATE TABLE collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    description TEXT NOT NULL DEFAULT '',
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_collections_user_id ON collections (user_id);

DROP TRIGGER IF EXISTS trg_collections_updated_at ON collections;
CREATE TRIGGER trg_collections_updated_at
    BEFORE UPDATE ON collections
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE collection_recipes (
    collection_id UUID NOT NULL REFERENCES collections (id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (collection_id, recipe_id)
);

CREATE INDEX idx_collection_recipes_recipe_id ON collection_recipes (recipe_id);

CREATE TABLE recipe_forks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_recipe_id UUID NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    child_recipe_id UUID NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (child_recipe_id),
    CHECK (parent_recipe_id != child_recipe_id)
);

CREATE INDEX idx_recipe_forks_parent ON recipe_forks (parent_recipe_id);
CREATE INDEX idx_recipe_forks_user_id ON recipe_forks (user_id);

CREATE TABLE recipe_transcripts (
    recipe_id UUID PRIMARY KEY REFERENCES recipes (id) ON DELETE CASCADE,
    segments JSONB NOT NULL DEFAULT '[]',
    raw_text TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL DEFAULT 'stream' CHECK (source IN ('stream', 'manual', 'ai')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS trg_recipe_transcripts_updated_at ON recipe_transcripts;
CREATE TRIGGER trg_recipe_transcripts_updated_at
    BEFORE UPDATE ON recipe_transcripts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
