alter table friendship
    add constraint friendship_unique unique (requested_from, requested_to);