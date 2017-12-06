
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE contexts RENAME COLUMN mid TO m_id;
ALTER TABLE contexts RENAME COLUMN qid TO q_id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE contexts RENAME COLUMN m_id TO mid;
ALTER TABLE contexts RENAME COLUMN q_id TO qid;
