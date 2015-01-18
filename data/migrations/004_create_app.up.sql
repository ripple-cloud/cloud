CREATE TABLE apps (
  id serial PRIMARY KEY NOT NULL,
  slug varchar(255) NOT NULL,
  hub_id int REFERENCES users(id) NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now()
);
CREATE UNIQUE INDEX index_apps_on_slug ON apps USING btree (slug);
CREATE UNIQUE INDEX index_apps_on_hub_id ON apps USING btree (hub_id);
