package storage

import (
	"RealWorldWeb/config"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *sqlx.DB //为什么这个会被引用？
var err error
var gormDB *gorm.DB

func init() {
	db, err = sqlx.Open("mysql", config.GetMysqlDSN()) //不加这个，时间解析有问题
	if err != nil {
		panic(fmt.Errorf("oepn mysql failed, error: %w", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping db failed, error: %w", err))
	}

	gormDB, err = gorm.Open(gorm_mysql.New(gorm_mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	gormDB = gormDB.Debug()
}

func IsNotFound(err error) bool {
	return errors.Is(err,
		gorm.ErrRecordNotFound)
}
