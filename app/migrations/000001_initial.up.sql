create table if not exists house (
    id serial primary key,
    address text not null,
    year integer not null,
    developer text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,

	unique (address)
);
