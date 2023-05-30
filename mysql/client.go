package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"tuiter.com/api/kit"
)

func Connect() *GormEngine {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/tuiter?parseTime=true"), &gorm.Config{})
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

func (g *GormEngine) MockData() error {
	absPath, _ := filepath.Abs("../tuiter-back/mysql/mock.sql")
	b, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	txExecution := g.gorm.Exec(string(b))
	return txExecution.Error
}

func (g *GormEngine) Offset(offset int) kit.DatabaseActions {
	g.gorm = g.gorm.Offset(offset)
	return g
}

func (g *GormEngine) Limit(limit int) kit.DatabaseActions {
	g.gorm = g.gorm.Limit(limit)
	return g
}
