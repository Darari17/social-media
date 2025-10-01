CREATE TABLE
  public.users (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    name character varying(100) NULL,
    email character varying(150) NOT NULL,
    password text NOT NULL,
    avatar text NULL,
    bio text NULL,
    created_at timestamp without time zone NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NULL
  );

ALTER TABLE
  public.users
ADD
  CONSTRAINT users_pkey PRIMARY KEY (id)