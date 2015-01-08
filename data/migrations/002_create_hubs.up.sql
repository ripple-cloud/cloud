CREATE TABLE hubs (
  id serial primary key,
  hub varchar(255) not null,
  user_id int references users(id)
);
