
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE claims ADD COLUMN image varchar(500);
ALTER TABLE users ADD COLUMN url varchar(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE claims DROP COLUMN image;
ALTER TABLE users DROP COLUMN url;