create table if not exists sdn (
    `uid` integer not null,
    `first_name` varchar(200),
    `last_name` varchar(200),
    unique key (`uid`)
);