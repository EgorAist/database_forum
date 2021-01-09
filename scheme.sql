
CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS forum_user;
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
  nickname citext PRIMARY KEY COLLATE "POSIX",
  fullname text,
  about text,
  email citext unique not null
);

CREATE INDEX idx_nick_nick ON "user" (nickname);
CREATE INDEX idx_nick_email ON "user" (email);
CREATE INDEX idx_nick_cover ON "user" (nickname, fullname, about, email);

CREATE TABLE forum (
  user_nick   citext references "user" not null,
  slug        citext PRIMARY KEY,
  title       text not null,
  thread_count integer default 0 not null,
  post_count integer default 0 not null
);

CREATE INDEX idx_forum_slug ON forum using hash(slug);

CREATE TABLE forum_user (
  nickname citext references "user",
  forum_slug citext references "forum",
  CONSTRAINT unique_forum_user UNIQUE (nickname, forum_slug)
);

CREATE INDEX idx_forum_user ON forum_user (nickname, forum_slug);

CREATE TABLE thread (
  id BIGSERIAL PRIMARY KEY,
  slug citext unique ,
  forum_slug citext references forum not null,
  user_nick citext references "user" not null,
  created timestamp with time zone default now(),
  title text not null,
  votes integer default 0 not null,
  message text not null
);

CREATE INDEX idx_thread_id ON thread(id);
CREATE INDEX idx_thread_slug ON thread(slug);
CREATE INDEX idx_thread_coverage ON thread (forum_slug, created, id, slug, user_nick, title, message, votes);

CREATE TABLE vote (
  nickname citext references "user",
  voice boolean not null,
  thread_id integer references thread,
  CONSTRAINT unique_vote UNIQUE (nickname, thread_id)
);

CREATE INDEX idx_vote ON vote(thread_id, voice);
/*
CREATE TABLE post (
  id BIGSERIAL PRIMARY KEY,
  path integer[],
  author citext references "user",
  created timestamp with time zone,
  edited boolean,
  message text,
  parent_id integer references post (id),
  forum_slug citext,
  thread_id integer references thread NOT NULL
);

CREATE INDEX ON post(thread_id, id, created, author, edited, message, parent_id, forum_slug);
CREATE INDEX idx_post_thread_id_p_i ON post(thread_id, (path[1]), id);
*/

CREATE TABLE post
(
    id        integer                            NOT NULL PRIMARY KEY,
    author citext references "user",
  --  author    text                               NOT NULL,
    created   text                               NOT NULL,
    forum_slug     citext                               NOT NULL,
    edited boolean   DEFAULT false            NOT NULL,
    message   text                               NOT NULL,
    parent_id    integer   DEFAULT 0                NOT NULL,
    thread_id    integer                            NOT NULL,
    path      INTEGER[] DEFAULT '{0}'::INTEGER[] NOT NULL
);

CREATE SEQUENCE post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE post_id_seq OWNED BY post.id;
ALTER TABLE ONLY post ALTER COLUMN id SET DEFAULT nextval('post_id_seq'::regclass);
SELECT pg_catalog.setval('post_id_seq', 1, false);

CREATE INDEX post_author_forum_index ON post USING btree (lower(author), lower(forum_slug));
CREATE INDEX post_forum_index ON post USING btree (lower(forum_slug));
CREATE INDEX post_parent_index ON post USING btree (parent_id);
CREATE INDEX post_path_index ON post USING gin (path);
CREATE INDEX post_thread_index ON post USING btree (thread_id);

CREATE OR REPLACE FUNCTION change_edited_post() RETURNS trigger as $change_edited_post$
BEGIN
  IF NEW.message <> OLD.message THEN
    NEW.edited = true;
  END IF;
  
  return NEW;
END;
$change_edited_post$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS change_edited_post ON post;

CREATE TRIGGER change_edited_post BEFORE UPDATE ON post
  FOR EACH ROW EXECUTE PROCEDURE change_edited_post();

/*CREATE OR REPLACE FUNCTION create_path() RETURNS trigger as $create_path$
BEGIN
   IF NEW.parent_id IS NULL THEN
     NEW.path := (ARRAY [NEW.id]);
     return NEW;
   end if;

   NEW.path := (SELECT array_append(p.path, NEW.id::integer)
                from post p where p.id = NEW.parent_id);
  RETURN NEW;
END;
$create_path$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS create_path ON post;

CREATE TRIGGER create_path BEFORE INSERT ON post
  FOR EACH ROW EXECUTE PROCEDURE create_path();
*/
CREATE TABLE post_count(count bigint);

CREATE FUNCTION post_count() RETURNS trigger
    LANGUAGE plpgsql AS
$$BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE post_count SET count = count + 1;

        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE post_count SET count = count - 1;

        RETURN OLD;
    ELSE
        UPDATE post_count SET count = 0;

        RETURN NULL;
    END IF;
END;$$;

CREATE CONSTRAINT TRIGGER post_count_mod
    AFTER INSERT OR DELETE ON post
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE PROCEDURE post_count();

-- TRUNCATE triggers must be FOR EACH STATEMENT
CREATE TRIGGER post_count_trunc AFTER TRUNCATE ON post
    FOR EACH STATEMENT EXECUTE PROCEDURE post_count();

-- initialize the counter table
INSERT INTO post_count
SELECT count(*) FROM post;