package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type queryBuilder string

func (q queryBuilder) orderBy(order string) queryBuilder {
	return queryBuilder(fmt.Sprintf("%s ORDER BY %s", q, order))
}

func (q queryBuilder) paginated(limit int, offset int) queryBuilder {
	return queryBuilder(fmt.Sprintf("%s LIMIT %d OFFSET %d;", q, limit, offset))
}

func Connect(user string, pass string, host string, db string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, db)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("failed to connect database: %v", err)
		panic("failed to connect database")
	}

	return database
}
