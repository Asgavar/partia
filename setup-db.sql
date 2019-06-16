CREATE TABLE member (
  id INTEGER PRIMARY KEY,
  password VARCHAR
);

CREATE TABLE member_isleader (
  member_id INTEGER REFERENCES member(id) PRIMARY KEY
);

CREATE TABLE member_lastactive (
  member_id INTEGER REFERENCES member(id) PRIMARY KEY,
  last_active TIMESTAMP
);

CREATE TABLE authority (
  id INTEGER PRIMARY KEY
);

CREATE TABLE project (
  id INTEGER PRIMARY KEY,
  authority INTEGER REFERENCES authority(id)
);

CREATE TYPE action_type AS ENUM ('support', 'protest');

CREATE TABLE action (
  id INTEGER PRIMARY KEY,
  proposed_by INTEGER REFERENCES member(id),
  project_id INTEGER REFERENCES project(id),
  of_type action_type
);

CREATE TABLE already_used_numbers (
  number INTEGER PRIMARY KEY
);

CREATE TABLE vote (
  member_id INTEGER REFERENCES member(id),
  action_id INTEGER REFERENCES action(id)
);

CREATE TABLE upvote () INHERITS(vote);
CREATE TABLE downvote () INHERITS(vote);

CREATE TABLE trolltracker (
  member_id INTEGER REFERENCES member(id),
  saldo INTEGER,
  ups INTEGER,
  downs INTEGER
);

CREATE OR REPLACE FUNCTION increment_trolltracker() RETURNS trigger AS $xD$
  DECLARE
  proposer INTEGER;
BEGIN
  proposer := (SELECT proposed_by FROM action WHERE action.id = NEW.action_id);
  UPDATE trolltracker SET saldo = saldo + 1 WHERE member_id = proposer;
  UPDATE trolltracker SET ups = ups + 1 WHERE member_id = proposer;
  RETURN NEW;
END
$xD$ LANGUAGE plpgsql;

CREATE TRIGGER upvote_trolltracker
  AFTER INSERT ON upvote
  FOR EACH ROW
    EXECUTE PROCEDURE increment_trolltracker();

CREATE OR REPLACE FUNCTION decrement_trolltracker() RETURNS trigger AS $xD$
  DECLARE
  proposer INTEGER;
BEGIN
  proposer := (SELECT proposed_by FROM action WHERE action.id = NEW.action_id);
  UPDATE trolltracker SET saldo = saldo - 1 WHERE member_id = proposer;
  UPDATE trolltracker SET downs = downs + 1 WHERE member_id = proposer;
  RETURN NEW;
END
$xD$ LANGUAGE plpgsql;

CREATE TRIGGER downvote_trolltracker
  AFTER INSERT ON downvote
  FOR EACH ROW
    EXECUTE PROCEDURE decrement_trolltracker();

CREATE OR REPLACE FUNCTION create_a_trolltracker_entry() RETURNS trigger AS $xD$
BEGIN
  INSERT INTO trolltracker VALUES (NEW.id, 0, 0, 0);
  RETURN NEW;
END
$xD$ LANGUAGE plpgsql;

CREATE TRIGGER member_create_a_trolltracker_entry
  AFTER INSERT ON member
  FOR EACH ROW
    EXECUTE PROCEDURE create_a_trolltracker_entry();
