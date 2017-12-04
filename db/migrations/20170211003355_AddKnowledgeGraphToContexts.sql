
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE contexts ADD COLUMN mid varchar(255);
ALTER TABLE contexts ADD COLUMN qid varchar(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE contexts DROP COLUMN mid;
ALTER TABLE contexts DROP COLUMN qid;
