create table article
(
    id              bigint(20)    not null primary key auto_increment,
    author_username varchar(255)  not null,
    title           varchar(4096) not null,
    slug            varchar(4096)  not null, /*这里扩容时会遇到：[42000][1071] Specified key was too long; max key length is 3072 bytes. 有2个方式：1.更换别的字符集，2.对这个列设置前缀索引*/
    `body`          text          not null,
    `description`   varchar(2048) not null default '',
    tag_list        varchar(1024) not null default '[]',
    created_at      datetime      not null default current_timestamp,
    updated_at      datetime      not null default current_timestamp
) engine = InnoDB
  default charset = utf8mb4;

