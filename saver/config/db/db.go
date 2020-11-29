package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func Get() *gorm.DB {
	return db
}

func Connect(url string) error {
	var err error
	db, err = gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	db.DB()
	err = db.DB().Ping()
	if err != nil {
		return err
	}
	db.LogMode(false)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)
	db.SingularTable(false)
	return nil
}

func Close() {
	db.Close()
}
