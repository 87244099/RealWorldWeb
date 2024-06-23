package models

import "time"

/**
create table article_comment
(
    id              bigint(20)    not null primary key auto_increment,
    author_username varchar(255)  not null,
    `body`          text          not null,
    created_at      datetime      not null default current_timestamp,
    updated_at      datetime      not null default current_timestamp
) engine = InnoDB
  default charset = utf8mb4;
*/

type ArticleComment struct {
	Id             int64  `gorm:"primaryKey"` //依赖db的字段自增
	AuthorUsername string `gorm:"column:author_username"`
	Body           string
	ArticleId      int64 `db:"article_id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	AuthorUserEmail string `gorm:"->"`
	AuthorUserImage string `gorm:"->"`
	AuthorUserBio   string `gorm:"->"`
}

func (a ArticleComment) TableName() string { //如果这里不指定的话，会默认为article_comments
	return "article_comment" //todo what is syntax sugar?
}
