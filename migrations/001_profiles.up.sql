create table if not exists profiles (
    id serial primary key,
    first_name varchar(128),
    last_name varchar(128),
    age integer,
    sex enum("male", "female"),
    about text,
    location varchar(100)
)