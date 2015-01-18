-- enable hstore extension
CREATE EXTENSION hstore;

CREATE TABLE received_messages (
  id serial PRIMARY KEY NOT NULL,
  hub_id int REFERENCES hubs(id) NOT NULL,
  topic varchar(255) NOT NULL,
  meta hstore,
  message text,
  received_at timestamp without time zone DEFAULT now()
);
