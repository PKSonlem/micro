create table if not exists flat (
    id serial primary key,
    house_id integer not null references house(id),
    flat_number integer not null,
    price integer not null,
    rooms integer not null,
    created_at timestamp with time zone not null,
    update_at timestamp with time zone not null,

    unique (house_id, flat_number)
);
