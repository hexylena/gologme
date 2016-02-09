package store

const DB_SCHEMA string = `
create table if not exists users (
    id integer not null primary key autoincrement,
    username text,
    api_key text
);

create table if not exists windowLogs (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    name text,
    foreign key (uid) references users(id)
);

create table if not exists keyLogs (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    count integer,
    foreign key (uid) references users(id)
);

create table if not exists notes (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    type integer,
    contents text,
    foreign key (uid) references users(id)
);
`
