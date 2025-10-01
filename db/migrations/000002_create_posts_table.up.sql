CREATE TABLE
  public.posts (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    user_id integer NOT NULL,
    content_text text NULL,
    content_image text NULL,
    created_at timestamp without time zone NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NULL,
    deleted_at timestamp without time zone NULL
  );

ALTER TABLE
  public.posts
ADD
  CONSTRAINT posts_pkey PRIMARY KEY (id)