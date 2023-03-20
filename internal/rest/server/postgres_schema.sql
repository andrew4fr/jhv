create table if not exists sdn (
    uid integer not null,
    first_name varchar(50),
    last_name varchar(50)
);

create unique index uid_idx on sdn(uid);