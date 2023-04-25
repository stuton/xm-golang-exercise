create table companies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(15) unique not null,
    description text,
    amount_employees integer not null,
    registered boolean not null,
    type varchar(15) not null
);

create table users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar(255) unique not null,
    password varchar(255) not null
);

insert into users (username, password) values ('admin', '$2a$14$86k0ddco8jPeyt0qUFhZ1.v0Ma4dMkGthjfy1h6qJTMOdDI/gnV8m');
