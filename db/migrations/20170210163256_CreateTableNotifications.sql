
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE notifications (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    user_id integer NOT NULL,
    type integer NOT NULL,
    item_id uuid,
    item_type integer,
    old_id uuid,
    old_type integer,
    new_id uuid,
    new_type integer,
    viewed boolean NOT NULL DEFAULT false
);

CREATE SEQUENCE notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE notifications_id_seq OWNED BY notifications.id;

ALTER TABLE ONLY notifications ALTER COLUMN id SET DEFAULT nextval('notifications_id_seq'::regclass);

SELECT pg_catalog.setval('notifications_id_seq', 1, false);

ALTER TABLE ONLY notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);

CREATE INDEX idx_notifications_user_id ON notifications USING btree (user_id);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE notifications;
