-- METADATA:
-- {
--   "description": "Setup script for types used in the media database",
--   "version": "1.0.0", 
--   "author": "artumont"
-- }

-- Table to store different media types
CREATE TABLE IF NOT EXISTS media_types(
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) UNIQUE NOT NULL -- e.g., 'movie', 'tv' or 'youtube'
);

-- Prepopulate media_types with common types
INSERT INTO media_types (type_name) VALUES 
    ('movie'),
    ('tv'),
    ('youtube')
ON CONFLICT (type_name) DO NOTHING;

-- Table to store genre IDs from TMDB
CREATE TABLE IF NOT EXISTS genre_ids(
    id INT UNIQUE NOT NULL, -- ID from TMDB
    genre_name VARCHAR(100) UNIQUE NOT NULL,
    media_type VARCHAR(50) REFERENCES media_types(type_name) ON DELETE CASCADE -- 'e.g., movie or tv' (not youtube as it has no genres)
);

-- Prepopulate genre_ids with common movie genres from TMDB
INSERT INTO genre_ids (id, media_type, genre_name) VALUES 
    (28, 'movie', 'Action'),
    (12, 'movie', 'Adventure'),
    (16, 'movie', 'Animation'),
    (35, 'movie', 'Comedy'),
    (80, 'movie', 'Crime'),
    (99, 'movie', 'Documentary'),
    (18, 'movie', 'Drama'),
    (10751, 'movie', 'Family'),
    (14, 'movie', 'Fantasy'),
    (36, 'movie', 'History'),
    (27, 'movie', 'Horror'),
    (10402, 'movie', 'Music'),
    (9648, 'movie', 'Mystery'),
    (10749, 'movie', 'Romance'),
    (878, 'movie', 'Science Fiction'),
    (10770, 'movie', 'TV Movie'),
    (53, 'movie', 'Thriller'),
    (10752, 'movie', 'War'),
    (37, 'movie', 'Western')
ON CONFLICT (id) DO NOTHING;

-- Prepopulate genre_ids with common TV genres from TMDB
INSERT INTO genre_ids (id, media_type, genre_name) VALUES 
    (10759, 'tv', 'Action & Adventure'),
    (16, 'tv', 'Animation'),
    (35, 'tv', 'Comedy'),
    (80, 'tv', 'Crime'),
    (99, 'tv', 'Documentary'),
    (18, 'tv', 'Drama'),
    (10751, 'tv', 'Family'),
    (10762, 'tv', 'Kids'),
    (9648, 'tv', 'Mystery'),
    (10763, 'tv', 'News'),
    (10764, 'tv', 'Reality'),
    (10765, 'tv', 'Sci-Fi & Fantasy'),
    (10766, 'tv', 'Soap'),
    (10767, 'tv', 'Talk'),
    (10768, 'tv', 'War & Politics'),
    (37, 'tv', 'Western')
ON CONFLICT (id) DO NOTHING;

-- Table to store supported languages
CREATE TABLE IF NOT EXISTS supported_languages(
    id SERIAL PRIMARY KEY,
    language_code VARCHAR(10) UNIQUE NOT NULL, -- e.g., 'en-US', 'es-ES' needs to match ISO 639-1 format
    language_name VARCHAR(100) NOT NULL -- e.g., 'English', 'Spanish'
);

-- Prepopulate supported_languages with some common languages
INSERT INTO supported_languages (language_code, language_name) VALUES 
    ('en-US', 'English'),
    ('es-ES', 'Spanish'),
    ('fr-FR', 'French'),
    ('de-DE', 'German'),
    ('it-IT', 'Italian'),
    ('ja-JP', 'Japanese'),
    ('zh-CN', 'Chinese'),
    ('ru-RU', 'Russian')
ON CONFLICT (language_code) DO NOTHING;