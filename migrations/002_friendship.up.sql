create table if not exists friendship(
    requested_from bigint unsigned not null references profiles(id),
    requested_to bigint unsigned not null references profiles(id),
    state enum("created", "accepted", "declined") default "created",
    created_at timestamp default current_timestamp,
    foreign key (requested_from) references profiles(id) on delete cascade,
    foreign key (requested_to) references profiles(id) on delete cascade
);