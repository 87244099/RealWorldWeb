package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

/**
create table article
(
    id              bigint(20)    not null primary key auto_increment,
    author_username varchar(255)  not null,
    title           varchar(1024) not null,
    slug            varchar(512)  not null unique,
    `body`          text          not null,
    `description`   varchar(2048) not null default '',
    tag_list        varchar(1024) not null default '[]',
    created_at      datetime      not null default current_timestamp,
    updated_at      datetime      not null default current_timestamp
) engine = InnoDB
  default charset = utf8mb4;
*/

type Article struct {
	Id             int64  `db:"id"` //依赖db的字段自增
	AuthorUsername string `gorm:"column:author_username"`
	Title          string
	Slug           string
	Body           string
	Description    string
	TagList        TagList `gorm:"type:string"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	AuthorUserEmail string `gorm:"->"`
	AuthorUserImage string `gorm:"->"`
	AuthorUserBio   string `gorm:"->"`
}

func (a Article) TableName() string {
	return "article" //todo what is syntax sugar?
}

type TagList []string

// from https://gorm.io/docs/data_types.html
// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *TagList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, j)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j TagList) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
