CREATE DATABASE IF NOT EXISTS auth CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

create table user
(
    id        int auto_increment,
    public_id varchar(255)                                     not null,
    email     varchar(255)                                     not null,
    role      enum ('admin', 'manager', 'user') default 'user' not null,
    constraint user_pk
        primary key (id),
    constraint user_public_id
        unique (public_id),
    constraint user_email
        unique (email)
);