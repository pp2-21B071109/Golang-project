ALTER TABLE coins ADD CONSTRAINT coins_coins_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE coins ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

CREATE INDEX IF NOT EXISTS coins_title_idx ON coins USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS coins_genres_idx ON coins USING GIN (genres);