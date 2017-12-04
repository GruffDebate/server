
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE claims ADD COLUMN truth_ru numeric;

ALTER TABLE arguments ADD COLUMN impact_ru numeric;
ALTER TABLE arguments ADD COLUMN relevance_ru numeric;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE claims DROP COLUMN truth_ru;

ALTER TABLE arguments DROP COLUMN impact_ru;
ALTER TABLE arguments DROP COLUMN relevance_ru;
