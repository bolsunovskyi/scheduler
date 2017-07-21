CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    email varchar(256) NOT NULL,
    first_name varchar(256) NULL,
    last_name varchar(256) NULL,
    password varchar(256) NOT NULL,
    active bool DEFAULT true,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_role role DEFAULT 'user'
);

CREATE TABLE session (
  id SERIAL PRIMARY KEY,
  value varchar(256) NOT NULL,
  ip varchar(256) NOT NULL,
  user_id INT NOT NULL,
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);

ALTER TABLE session ADD CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE;