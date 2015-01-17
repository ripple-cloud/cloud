CREATE TABLE hubs (
  id serial PRIMARY KEY NOT NULL,
  slug varchar(255) NOT NULL,
  user_id int REFERENCES users(id) NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now()
);
