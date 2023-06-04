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
	Create(value interface{}) databaseActions
}

type Reader interface {
	First(dest interface{}, conds ...interface{}) databaseActions
	Find(dest interface{}, conds ...interface{}) databaseActions
	Search(dest interface{}, query map[string]interface{}) databaseActions
	Offset(offset int) databaseActions
	Limit(limit int) databaseActions
}

type databaseActions interface {
	Creator
	Reader
	Error() error
	AutoMigrate(dst ...interface{}) error
}

type GormEngine struct {
	gorm *gorm.DB
}

func (g *GormEngine) Error() error {
	return g.gorm.Error
}

func (g *GormEngine) Create(value interface{}) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.Create(value)

	return g
}

func (g *GormEngine) First(dest interface{}, conds ...interface{}) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.First(dest, conds...)

	return g
}

func (g *GormEngine) AutoMigrate(dst ...interface{}) error {
	err := g.gorm.AutoMigrate(dst...)

	if err != nil {
		return fmt.Errorf("error migrating %w", err)
	}

	return nil
}

func (g *GormEngine) Find(dest interface{}, conds ...interface{}) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.Where(conds).Find(dest)

	return g
}

func (g *GormEngine) Search(dest interface{}, query map[string]interface{}) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.Where(query).Find(dest)

	return g
}

func (g *GormEngine) MockData() error {
	absPath, _ := filepath.Abs("../tuiter-back/mysql/mock.sql")
	sqlScriptBytes, err := os.ReadFile(absPath)

	if err != nil {
		return fmt.Errorf("error reading mock file %w", err)
	}

	txExecution := g.gorm.Exec(string(sqlScriptBytes))

	if txExecution.Error != nil {
		return fmt.Errorf("error executing mock file %w", txExecution.Error)
	}

	return nil
}

func (g *GormEngine) Offset(offset int) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.Offset(offset)

	return g
}

func (g *GormEngine) Limit(limit int) databaseActions { //nolint: ireturn
	g.gorm = g.gorm.Limit(limit)

	return g
}