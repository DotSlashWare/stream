-- METADATA:
-- {
--   "description": "Setup script for the media server profiles",
--   "version": "1.0.0", 
--   "author": "artumont",
--   "dependencies": ["01_type_setup.sql"]
-- }

-- Table to store all user profiles
CREATE TABLE IF NOT EXISTS profiles(
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    totp_secret VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login TIMESTAMPTZ
);

-- Table to store profile settings and preferences
CREATE TABLE IF NOT EXISTS profile_settings(
    profile_id INT PRIMARY KEY REFERENCES Profiles(id) ON DELETE CASCADE,
    is_over_18 BOOLEAN NOT NULL DEFAULT FALSE,
    preferred_media_type INT REFERENCES media_types(id),
    preferred_language INT REFERENCES supported_languages(id),
    autoplay_next_episode BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(profile_id)
);

-- Table to store watch history for profiles
CREATE TABLE IF NOT EXISTS profile_watch_history(
    id BIGSERIAL PRIMARY KEY,
    profile_id INT REFERENCES Profiles(id) ON DELETE CASCADE,
    media_id VARCHAR(50) NOT NULL,
    media_type INT REFERENCES media_types(id),
    media_progress INT NOT NULL DEFAULT 0, -- in seconds
    media_cache JSONB, -- cache metadata to avoid extra fetches (stuff like title, poster, etc.)
    watched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(profile_id, media_id)
);