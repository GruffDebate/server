-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Table Users
--
CREATE TABLE users (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    name character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    hashed_password bytea NOT NULL,
    image character varying(255),
    curator boolean NOT NULL DEFAULT false,
    admin boolean NOT NULL DEFAULT false,
    email_verified_at timestamp with time zone
);

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

SELECT pg_catalog.setval('users_id_seq', 1, false);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX uix_users_email ON users USING btree (email);
CREATE UNIQUE INDEX uix_users_username ON users USING btree (username);

--
-- Table Debates
--
CREATE TABLE debates (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    created_by_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(1000) NOT NULL,
    description character varying(4000),
    truth numeric
);

ALTER TABLE ONLY debates
    ADD CONSTRAINT debates_pkey PRIMARY KEY (id);

ALTER TABLE ONLY debates
    ADD CONSTRAINT debates_created_by_id_fkey FOREIGN KEY (created_by_id)
      REFERENCES users (id);

CREATE INDEX idx_debates_created_by_id ON debates USING btree (created_by_id);

-- 
-- Table Arguments
-- 
CREATE TABLE arguments (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    created_by_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(1000) NOT NULL,
    description character varying(4000),
    type integer NOT NULL,
    relevance numeric,
    impact numeric,
    parent_id uuid,
    argument_id uuid,
    debate_id uuid NOT NULL
);

ALTER TABLE ONLY arguments
    ADD CONSTRAINT arguments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY arguments
    ADD CONSTRAINT arguments_parent_id_fkey FOREIGN KEY (parent_id)
      REFERENCES debates (id);

ALTER TABLE ONLY arguments
    ADD CONSTRAINT arguments_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

ALTER TABLE ONLY arguments
    ADD CONSTRAINT arguments_argument_id_fkey FOREIGN KEY (argument_id)
      REFERENCES arguments (id);

ALTER TABLE ONLY arguments
    ADD CONSTRAINT arguments_created_by_id_fkey FOREIGN KEY (created_by_id)
      REFERENCES users (id);

CREATE INDEX idx_arguments_parent_id ON arguments USING btree (parent_id);
CREATE INDEX idx_arguments_debate_id ON arguments USING btree (debate_id);
CREATE INDEX idx_arguments_argument_id ON arguments USING btree (argument_id);
CREATE INDEX idx_arguments_created_by_id ON arguments USING btree (created_by_id);

CREATE UNIQUE INDEX uix_arguments_type ON arguments USING btree (parent_id, debate_id, argument_id, type);

-- 
-- Table debate_opinions
-- 
CREATE TABLE debate_opinions (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    truth numeric,
    user_id integer NOT NULL,
    debate_id uuid NOT NULL
);

CREATE SEQUENCE debate_opinions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE debate_opinions_id_seq OWNED BY debate_opinions.id;

ALTER TABLE ONLY debate_opinions ALTER COLUMN id SET DEFAULT nextval('debate_opinions_id_seq'::regclass);

SELECT pg_catalog.setval('debate_opinions_id_seq', 1, false);

ALTER TABLE ONLY debate_opinions
    ADD CONSTRAINT debate_opinions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY debate_opinions
    ADD CONSTRAINT debate_opinionss_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES users (id);

ALTER TABLE ONLY debate_opinions
    ADD CONSTRAINT debate_opinions_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

CREATE INDEX idx_debate_opinions_parent_id ON debate_opinions USING btree (user_id);
CREATE INDEX idx_debate_opinions_debate_id ON debate_opinions USING btree (debate_id);

-- 
-- Table argument_opinions
-- 
CREATE TABLE argument_opinions (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    relevance numeric,
    impact numeric,
    user_id integer NOT NULL,
    argument_id uuid NOT NULL
);

CREATE SEQUENCE argument_opinions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE argument_opinions_id_seq OWNED BY argument_opinions.id;

ALTER TABLE ONLY argument_opinions ALTER COLUMN id SET DEFAULT nextval('argument_opinions_id_seq'::regclass);

SELECT pg_catalog.setval('argument_opinions_id_seq', 1, false);

ALTER TABLE ONLY argument_opinions
    ADD CONSTRAINT argument_opinions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY argument_opinions
    ADD CONSTRAINT argument_opinionss_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES users (id);

ALTER TABLE ONLY argument_opinions
    ADD CONSTRAINT argument_opinions_argument_id_fkey FOREIGN KEY (argument_id)
      REFERENCES arguments (id);

CREATE INDEX idx_argument_opinions_parent_id ON argument_opinions USING btree (user_id);
CREATE INDEX idx_argument_opinions_argument_id ON argument_opinions USING btree (argument_id);

-- 
-- Table Links
-- 
CREATE TABLE links (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(1000) NOT NULL,
    description character varying(4000),
    url character varying(4000),
    created_by_id integer NOT NULL,
    debate_id uuid NOT NULL
);

ALTER TABLE ONLY links
    ADD CONSTRAINT links_pkey PRIMARY KEY (id);

ALTER TABLE ONLY links
    ADD CONSTRAINT links_created_by_id_fkey FOREIGN KEY (created_by_id)
      REFERENCES users (id);

ALTER TABLE ONLY links
    ADD CONSTRAINT links_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

CREATE INDEX idx_links_created_by_id ON links USING btree (created_by_id);
CREATE INDEX idx_links_debate_id ON links USING btree (debate_id);

-- 
-- Table Tags
-- 
CREATE TABLE tags (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(50) NOT NULL
);

CREATE SEQUENCE tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tags_id_seq OWNED BY tags.id;

ALTER TABLE ONLY tags ALTER COLUMN id SET DEFAULT nextval('tags_id_seq'::regclass);

SELECT pg_catalog.setval('tags_id_seq', 1, false);

ALTER TABLE ONLY tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


-- 
-- Table Tags
-- 
CREATE TABLE contexts (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(50) NOT NULL,
    description character varying(4000),
    parent_id integer
);

CREATE SEQUENCE contexts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contexts_id_seq OWNED BY contexts.id;

ALTER TABLE ONLY contexts ALTER COLUMN id SET DEFAULT nextval('contexts_id_seq'::regclass);

SELECT pg_catalog.setval('contexts_id_seq', 1, false);

ALTER TABLE ONLY contexts
    ADD CONSTRAINT contexts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY contexts
    ADD CONSTRAINT contexts_parent_id_fkey FOREIGN KEY (parent_id)
      REFERENCES contexts (id);

CREATE INDEX idx_contexts_parent_id ON contexts USING btree (parent_id);

-- 
-- Table Values
-- 
CREATE TABLE "values" (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    title character varying(1000) NOT NULL,
    description character varying(4000),
    parent_id integer
);

CREATE SEQUENCE values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE values_id_seq OWNED BY "values".id;

ALTER TABLE ONLY "values" ALTER COLUMN id SET DEFAULT nextval('values_id_seq'::regclass);

SELECT pg_catalog.setval('values_id_seq', 1, false);

ALTER TABLE ONLY "values"
    ADD CONSTRAINT values_pkey PRIMARY KEY (id);

ALTER TABLE ONLY "values"
    ADD CONSTRAINT values_parent_id_fkey FOREIGN KEY (parent_id)
      REFERENCES "values" (id);

CREATE INDEX idx_values_parent_id ON values USING btree (parent_id);


-- 
-- Table debate_tags
-- 
CREATE TABLE debate_tags (
    tag_id integer NOT NULL,
    debate_id uuid NOT NULL
);

ALTER TABLE ONLY debate_tags
    ADD CONSTRAINT debate_tags_pkey PRIMARY KEY (tag_id, debate_id);

ALTER TABLE ONLY debate_tags
    ADD CONSTRAINT debate_tags_tag_id_fkey FOREIGN KEY (tag_id)
      REFERENCES tags (id);

ALTER TABLE ONLY debate_tags
    ADD CONSTRAINT debate_tags_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

CREATE INDEX idx_debate_tags_tag_id ON debate_tags USING btree (tag_id);
CREATE INDEX idx_debate_tags_debate_id ON debate_tags USING btree (debate_id);

-- 
-- Table debate_contexts
-- 
CREATE TABLE debate_contexts (
    context_id integer NOT NULL,
    debate_id uuid NOT NULL
);

ALTER TABLE ONLY debate_contexts
    ADD CONSTRAINT debate_contexts_pkey PRIMARY KEY (context_id, debate_id);

ALTER TABLE ONLY debate_contexts
    ADD CONSTRAINT debate_contexts_context_id_fkey FOREIGN KEY (context_id)
      REFERENCES contexts (id);

ALTER TABLE ONLY debate_contexts
    ADD CONSTRAINT debate_contexts_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

CREATE INDEX idx_debate_contexts_context_id ON debate_contexts USING btree (context_id);
CREATE INDEX idx_debate_contexts_debate_id ON debate_contexts USING btree (debate_id);

-- 
-- Table debate_values
-- 
CREATE TABLE debate_values (
    value_id integer NOT NULL,
    debate_id uuid NOT NULL
);

ALTER TABLE ONLY debate_values
    ADD CONSTRAINT debate_values_pkey PRIMARY KEY (value_id, debate_id);

ALTER TABLE ONLY debate_values
    ADD CONSTRAINT debate_values_value_id_fkey FOREIGN KEY (value_id)
      REFERENCES values (id);

ALTER TABLE ONLY debate_values
    ADD CONSTRAINT debate_values_debate_id_fkey FOREIGN KEY (debate_id)
      REFERENCES debates (id);

CREATE INDEX idx_debate_values_value_id ON debate_values USING btree (value_id);
CREATE INDEX idx_debate_values_debate_id ON debate_values USING btree (debate_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE debate_values;
DROP TABLE debate_contexts;
DROP TABLE debate_tags;
DROP TABLE "values";
DROP TABLE contexts;
DROP TABLE tags;
DROP TABLE links;
DROP TABLE argument_opinions;
DROP TABLE debate_opinions;
DROP TABLE arguments;
DROP TABLE debates;
DROP TABLE users;
