CREATE TABLE
  public.likes (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    user_id integer NOT NULL,
    post_id integer NOT NULL,
    created_at timestamp without time zone NULL DEFAULT CURRENT_TIMESTAMP
  );

ALTER TABLE
  public.likes
ADD
  CONSTRAINT likes_pkey PRIMARY KEY (id)