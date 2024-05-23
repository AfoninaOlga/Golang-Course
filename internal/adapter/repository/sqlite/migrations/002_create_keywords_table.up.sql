create table if not exists Keywords
(
    id integer primary key autoincrement,
    word text unique not null
);