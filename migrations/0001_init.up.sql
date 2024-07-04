create table if not exists person
(
    id serial primary key,
    surname varchar(255) not null,
    name varchar(255) not null,
    patronymic varchar(255),
    address text not null,
    passport_number varchar(11) unique
);

create table if not exists task
(
    id serial primary key,
    name text not null,
    start_tracking timestamp,
    stop_tracking timestamp,
    user_id bigint not null references person(id)
);