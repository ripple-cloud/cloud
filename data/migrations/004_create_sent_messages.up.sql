CREATE TABLE sent_messages (
  id serial PRIMARY KEY NOT NULL,
  topic varchar(255) NOT NULL,
  message text,
  hub_id int REFERENCES hubs(id),
  user_id int REFERENCES users(id) NOT NULL,
  requested_at timestamp without time zone DEFAULT now(),
);
