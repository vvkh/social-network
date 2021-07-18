create table if not exists users (
    id serial primary key,
    username varchar(100) unique,
    password binary(60) -- 60 is max size for bcrypt hash
);

alter table profiles
    add user_id bigint unsigned not null; -- it's safe at this point as we don't have any existing profiles anyway

alter table profiles
    add constraint profile_user_fk foreign key (user_id) references users(id) on delete cascade;