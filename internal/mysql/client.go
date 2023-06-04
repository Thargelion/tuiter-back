package mysql

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *GormEngine {
	database, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/tuiter?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err != nil {
		panic("failed to connect migrate")
	}

	return &GormEngine{database}
}

type Creator interface {
	Create(value interface{}) *GormEngine
}

type Reader interface {
	First(dest interface{}, conds ...interface{}) *GormEngine
	Find(dest interface{}, conds ...interface{}) *GormEngine
	Search(dest interface{}, query map[string]interface{}) *GormEngine
	Offset(offset int) *GormEngine
	Limit(limit int) *GormEngine
}

type Crud interface {
	Creator
	Reader
}

type DatabaseActions interface {
	Crud
	Error() error
	AutoMigrate(dst ...interface{}) error
}

type GormEngine struct {
	Gorm *gorm.DB
}

func (g *GormEngine) Error() error {
	return g.Gorm.Error
}

func (g *GormEngine) Create(value interface{}) *GormEngine {
	g.Gorm = g.Gorm.Create(value)

	return g
}

func (g *GormEngine) First(dest interface{}, conds ...interface{}) *GormEngine {
	g.Gorm = g.Gorm.First(dest, conds...)

	return g
}

func (g *GormEngine) AutoMigrate(dst ...interface{}) error {
	err := g.Gorm.AutoMigrate(dst...)

	if err != nil {
		return fmt.Errorf("error migrating %w", err)
	}

	return nil
}

func (g *GormEngine) Find(dest interface{}, conds ...interface{}) *GormEngine {
	g.Gorm = g.Gorm.Where(conds).Find(dest)

	return g
}

func (g *GormEngine) Search(dest interface{}, query map[string]interface{}) *GormEngine {
	g.Gorm = g.Gorm.Where(query).Find(dest)

	return g
}

func (g *GormEngine) MockData() error {
	absPath, _ := filepath.Abs("../tuiter-back/mysql/mock.sql")
	sqlScriptBytes, err := os.ReadFile(absPath)

	if err != nil {
		return fmt.Errorf("error reading mock file %w", err)
	}

	txExecution := g.Gorm.Exec(string(sqlScriptBytes))

	if txExecution.Error != nil {
		return fmt.Errorf("error executing mock file %w", txExecution.Error)
	}

	return nil
}

func (g *GormEngine) Offset(offset int) *GormEngine {
	g.Gorm = g.Gorm.Offset(offset)

	return g
}

func (g *GormEngine) Limit(limit int) *GormEngine {
	g.Gorm = g.Gorm.Limit(limit)

	return g
}
