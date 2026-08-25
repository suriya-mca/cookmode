CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    cuisine TEXT NOT NULL DEFAULT '',
    prep_time_min INT NOT NULL DEFAULT 0,
    cook_time_min INT NOT NULL DEFAULT 0,
    servings INT NOT NULL DEFAULT 1,
    difficulty TEXT NOT NULL DEFAULT '',
    dietary_tags TEXT[] NOT NULL DEFAULT '{}',
    ingredients JSONB NOT NULL DEFAULT '[]',
    steps JSONB NOT NULL DEFAULT '[]',
    substitutions JSONB NOT NULL DEFAULT '[]',
    nutrition JSONB,
    video_uid TEXT NOT NULL DEFAULT '',
    video_hls_url TEXT NOT NULL DEFAULT '',
    video_thumbnail_url TEXT NOT NULL DEFAULT '',
    video_duration_sec DOUBLE PRECISION NOT NULL DEFAULT 0,
    views INT NOT NULL DEFAULT 0,
    saves INT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'processing', 'published', 'archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_recipes_user_id ON recipes (user_id);
CREATE INDEX idx_recipes_cuisine ON recipes (cuisine);
CREATE INDEX idx_recipes_status ON recipes (status);
CREATE INDEX idx_recipes_dietary_tags ON recipes USING GIN (dietary_tags);

-- "I made this" photo posts
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    recipe_id UUID NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    rating_value INT CHECK (rating_value BETWEEN 0 AND 5),
    status TEXT NOT NULL DEFAULT 'visible' CHECK (status IN ('visible', 'archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_posts_recipe_id ON posts (recipe_id);
CREATE INDEX idx_posts_user_id ON posts (user_id);

-- Shopping lists (items stored as JSONB for flexibility)
CREATE TABLE shopping_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    items JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shopping_lists_user_id ON shopping_lists (user_id);
