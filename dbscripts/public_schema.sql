SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET client_min_messages = warning;
SET row_security = off;
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;
SET search_path = public, pg_catalog;
SET default_tablespace = '';

-- runners
CREATE TABLE runners (
  -- uuid_generate_vlmc() creates a new UUID during INSERT cmd.
  id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
  first_name text NOT NULL,
  last_name text NOT NULL,
  age integer,
  is_active boolean DEFAULT TRUE,
  country text NOT NULL,
  personal_best interval,
  season_best interval,
  -- set primary key for runners table to the id
  CONSTRAINT runners_pk PRIMARY KEY (id)
);
-- create indices on country and season_best for search for runner by country and season_best.
CREATE INDEX runners_country ON runners (country);
CREATE INDEX runners_season_best ON runners (season_best);

-- results
CREATE TABLE results (
  id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
  runner_id uuid NOT NULL,
  race_result interval NOT NULL,
  location text NOT NULL,
  position integer,
  year integer NOT NULL,
  -- set primary key for results table to id. Set foreign key for runner_id col.
  CONSTRAINT results_pk PRIMARY KEY (id),
  CONSTRAINT fk_results_runner_id FOREIGN KEY (runner_id)
    REFERENCES runners (id) MATCH SIMPLE
    -- define what happens if the row referenced by a foreign key is deleted. 
    -- no need for update because runners are never deleted just change is_active to false.
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);