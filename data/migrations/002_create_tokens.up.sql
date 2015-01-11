CREATE TABLE tokens (
  id bigserial PRIMARY KEY NOT NULL,
  user_id bigint REFERENCES users(id) NOT NULL,
  expires_in bigint,
  created_at timestamp without time zone DEFAULT now(),
  revoked_at timestamp without time zone
);
