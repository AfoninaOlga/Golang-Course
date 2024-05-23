create table if not exists ComicsKeywords
(
    id integer primary key autoincrement,
    comic_id integer,
    keyword_id integer,
    foreign key (comic_id) references Comics (id) on delete cascade,
    foreign key (keyword_id) references Keywords (id) on delete cascade
);