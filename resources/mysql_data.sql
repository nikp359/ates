CREATE DATABASE IF NOT EXISTS auth CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

create table if not exists user
(
    id          int auto_increment primary key,
    public_id   varchar(255)                                     not null,
    email       varchar(255)                                     not null,
    role        enum ('admin', 'manager', 'employee') default 'employee' not null,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    unique user_public_id (public_id),
    unique user_email (email),
    index user_updated_at (updated_at)
);

CREATE DATABASE IF NOT EXISTS task CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

create table if not exists user
(
    id          int auto_increment primary key,
    public_id   varchar(255)                                    not null,
    email       varchar(255)                                    not null,
    role        varchar(255)                                    not null,
    updated_at  TIMESTAMP not null ,
    unique user_public_id (public_id),
    unique user_email (email),
    index user_updated_at (updated_at)
);

create table if not exists task
(
  id                int auto_increment primary key,
  public_id         varchar(255)        not null,
  title             varchar(255)        not null,
  jira_id           varchar(255),
  description       varchar(255),
  status            enum ('new', 'completed') default 'new' not null,
  assigned_user_id  varchar(255)        not null,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  unique task_public_id (public_id),
  index task_assigned_user_id (assigned_user_id),
  index task_updated_at (updated_at)
);