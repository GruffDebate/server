
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE arguments DROP COLUMN relevance;
ALTER TABLE arguments DROP COLUMN relevance_ru;
ALTER TABLE arguments RENAME COLUMN impact TO strength;
ALTER TABLE arguments RENAME COLUMN impact_ru TO strength_ru;

ALTER TABLE argument_opinions DROP COLUMN relevance;
ALTER TABLE argument_opinions RENAME COLUMN impact TO strength;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE arguments ADD COLUMN relevance numeric;
ALTER TABLE arguments ADD COLUMN relevance_ru numeric;
ALTER TABLE arguments RENAME COLUMN strength TO impact;
ALTER TABLE arguments RENAME COLUMN strength_ru TO impact_ru;

ALTER TABLE argument_opinions ADD COLUMN relevance numeric;
ALTER TABLE argument_opinions RENAME COLUMN strength TO impact;
