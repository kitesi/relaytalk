create table server (
    id serial primary key,
    name text not null,
    description text,
    created_at timestamp with time zone default current_timestamp,
    owner_id integer not null references users(id),
    unique (name, owner_id) -- Ensures that a user cannot create multiple servers with the same name
);

create table channels (
    id serial primary key,
    name text unique not null,
    description text,
    created_at timestamp with time zone default current_timestamp,
    server_id integer not null references server(id),
    owner_id integer not null references users(id)
);

create table messages (
    id serial primary key,
    user_id integer not null references users(id),
    content text not null,
    channel_id integer not null references channels(id)
);
