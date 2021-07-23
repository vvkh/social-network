alter table friendship
    add constraint if not exists friendship_unique unique (requested_from, requested_to);