package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/tuiter"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err != nil {
		panic("failed to connect migrate")
	}

	return db
}
