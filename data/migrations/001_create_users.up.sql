CREATE TABLE users (
  id bigserial PRIMARY KEY NOT NULL,
  username varchar(255) NOT NULL UNIQUE,
  email varchar(255) NOT NULL UNIQUE,
  encrypted_password varchar(255) NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now()
);
CREATE UNIQUE INDEX index_users_on_username ON users USING btree (username);
CREATE UNIQUE INDEX index_users_on_email ON users USING btree (email);
