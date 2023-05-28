package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tuiter.com/api/kit"
)

func Connect() kit.Dao {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/tuiter"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err != nil {
		panic("failed to connect migrate")
	}

	return &GormEngine{db}
}

type GormEngine struct {
	gorm *gorm.DB
}

func (g *GormEngine) Error() error {
	return g.gorm.Error
}

func (g *GormEngine) Create(value interface{}) kit.Dao {
	g.gorm = g.gorm.Create(value)
	return g
}

func (g *GormEngine) First(dest interface{}, conds ...interface{}) kit.Dao {
	g.gorm = g.gorm.First(dest, conds...)
	return g
}

func (g *GormEngine) AutoMigrate(dst ...interface{}) error {
	return g.gorm.AutoMigrate(dst...)
}
