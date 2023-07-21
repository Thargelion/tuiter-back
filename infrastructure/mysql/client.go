package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	database, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/tuiter?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err != nil {
		panic("failed to connect migrate")
	}

	return database
}
