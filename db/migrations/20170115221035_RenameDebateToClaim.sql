
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE arguments RENAME COLUMN debate_id TO claim_id;
ALTER TABLE arguments RENAME COLUMN target_debate_id TO target_claim_id;
ALTER TABLE links RENAME COLUMN debate_id TO claim_id;
ALTER TABLE debate_opinions RENAME COLUMN debate_id TO claim_id;
ALTER TABLE debate_contexts RENAME COLUMN debate_id TO claim_id;
ALTER TABLE debate_tags RENAME COLUMN debate_id TO claim_id;
ALTER TABLE debate_values RENAME COLUMN debate_id TO claim_id;

ALTER TABLE debates RENAME TO claims;
ALTER TABLE debate_opinions RENAME TO claim_opinions;
ALTER TABLE debate_contexts RENAME TO claim_contexts;
ALTER TABLE debate_tags RENAME TO claim_tags;
ALTER TABLE debate_values RENAME TO claim_values;
ALTER SEQUENCE debate_opinions_id_seq RENAME TO claim_opinions_id_seq;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE claims RENAME TO debates;
ALTER TABLE claim_opinions RENAME TO debate_opinions;
ALTER TABLE claim_contexts RENAME TO debate_contexts;
ALTER TABLE claim_tags RENAME TO debate_tags;
ALTER TABLE claim_values RENAME TO debate_values;
ALTER SEQUENCE claim_opinions_id_seq RENAME TO debate_opinions_id_seq;

ALTER TABLE arguments RENAME COLUMN claim_id TO debate_id;
ALTER TABLE arguments RENAME COLUMN target_claim_id TO target_debate_id;
ALTER TABLE links RENAME COLUMN claim_id TO debate_id;
ALTER TABLE debate_opinions RENAME COLUMN claim_id TO debate_id;
ALTER TABLE debate_contexts RENAME COLUMN claim_id TO debate_id;
ALTER TABLE debate_tags RENAME COLUMN claim_id TO debate_id;
ALTER TABLE debate_values RENAME COLUMN claim_id TO debate_id;

