create table users (
    id serial primary key,
    username text unique not null,
    password_hash text not null,
    email text unique not null,
    email_verified boolean default false
)
