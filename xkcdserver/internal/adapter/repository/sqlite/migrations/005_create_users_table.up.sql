create table if not exists Users (
    name text primary key,
    password text not null,
    is_admin integer not null check ( is_admin in (0, 1) )
);