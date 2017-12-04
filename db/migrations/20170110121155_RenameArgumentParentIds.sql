
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE arguments RENAME COLUMN parent_id TO target_debate_id;
ALTER TABLE arguments RENAME COLUMN argument_id TO target_argument_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE arguments RENAME COLUMN target_debate_id TO parent_id;
ALTER TABLE arguments RENAME COLUMN target_argument_id TO argument_id;
