
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE contexts ADD COLUMN url varchar(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE contexts DROP COLUMN url;
