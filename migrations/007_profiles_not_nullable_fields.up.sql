alter table profiles modify first_name varchar(128) not null;
alter table profiles modify last_name varchar(128) not null;
alter table profiles modify age int not null;
alter table profiles modify sex enum("male", "female") not null;
