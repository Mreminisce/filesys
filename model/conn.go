package model

import (
	"fmt"
	"time"

	"filesys/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", config.MySQLDSN)
	if err != nil {
		fmt.Println("Connect database  failed: ", err)
		panic("Connect database  failed...")
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.LogMode(true)
	db = db.Set("gorm:table_options", "ENGINE=InnoDB;").AutoMigrate()
	DB = db
	return db
}
