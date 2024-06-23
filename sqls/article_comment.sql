create table article_comment
(
    id              bigint(20)    not null primary key auto_increment,
    author_username varchar(255)  not null,
    article_id      bigint        not null,
    `body`          text          not null,
    created_at      datetime      not null default current_timestamp,
    updated_at      datetime      not null default current_timestamp
) engine = InnoDB
  default charset = utf8mb4;