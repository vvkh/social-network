create table if not exists chats(
    id serial primary key,
    title varchar(255)
);

create table if not exists chat_members(
    chat_id bigint references chats(id),
    member_profile_id bigint references profiles(id),
    constraint chat_member_unique unique (chat_id, member_profile_id)
);

create table if not exists messages(
    id serial primary key,
    chat_id bigint,
    author_profile_id bigint,
    content text,
    was_read boolean,
    sent_at timestamp default now()
);