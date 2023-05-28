package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tuiter.com/api/kit"
)

func Connect() kit.DatabaseActions {
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

func (g *GormEngine) Create(value interface{}) kit.DatabaseActions {
	g.gorm = g.gorm.Create(value)
	return g
}

func (g *GormEngine) First(dest interface{}, conds ...interface{}) kit.DatabaseActions {
	g.gorm = g.gorm.First(dest, conds...)
	return g
}

func (g *GormEngine) AutoMigrate(dst ...interface{}) error {
	return g.gorm.AutoMigrate(dst...)
}

func (g *GormEngine) Find(dest interface{}, conds ...interface{}) kit.DatabaseActions {
	g.gorm = g.gorm.Where(conds).Find(dest)
	return g
}

func (g *GormEngine) Search(dest interface{}, query map[string]interface{}) kit.DatabaseActions {
	g.gorm = g.gorm.Where(query).Find(dest)
	return g
}
