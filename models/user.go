package models

/**
create table user (
      username varchar(64) not null primary key,
      password varchar(512) not null default '',
      email varchar(256) not null default '' unique,
      image varchar(1024) not null default ''
      bio varchar(1024) not null default ''
      ) engine=InnoDB default charset=utf8mb4
*/

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}
