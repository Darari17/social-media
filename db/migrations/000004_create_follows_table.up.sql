CREATE TABLE
  public.follows (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    follower_id integer NOT NULL,
    following_id integer NOT NULL,
    created_at timestamp without time zone NULL DEFAULT CURRENT_TIMESTAMP
  );

ALTER TABLE
  public.follows
ADD
  CONSTRAINT follows_pkey PRIMARY KEY (id)