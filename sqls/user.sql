create table user (
                      username varchar(64) not null primary key,
                      password varchar(512) not null default '',
                      email varchar(256) not null default '' unique,
                      image varchar(1024) not null default '',
                          bio varchar(1024) not null default ''
) engine=InnoDB default charset=utf8mb4
