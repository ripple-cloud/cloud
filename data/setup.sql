create table users (
  id serial primary key,
  username varchar(255) not null unique, 
  email varchar(255) not null unique,
  password varchar(255) not null,
  token varchar(255),
  created_at timestamp not null
);
