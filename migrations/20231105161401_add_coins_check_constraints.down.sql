ALTER TABLE coins DROP CONSTRAINT IF EXISTS coins_year_check;
ALTER TABLE coins DROP CONSTRAINT IF EXISTS coins_length_check;
DROP INDEX IF EXISTS movies_title_idx;
DROP INDEX IF EXISTS movies_genres_idx;