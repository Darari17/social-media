CREATE TABLE
  public.comments (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    user_id integer NOT NULL,
    post_id integer NOT NULL,
    content text NOT NULL,
    created_at timestamp without time zone NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NULL
  );

ALTER TABLE
  public.comments
ADD
  CONSTRAINT comments_pkey PRIMARY KEY (id)