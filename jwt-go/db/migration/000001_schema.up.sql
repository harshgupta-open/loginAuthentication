create TABLE IF NOT EXISTS "user_details" (
    id serial  primary key ,
    email varchar not null UNIQUE,
    user_password varchar,
    created date,
    updated date,
    deleted date
);