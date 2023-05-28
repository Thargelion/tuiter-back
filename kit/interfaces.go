package kit

type Creator interface {
	Create(value interface{}) Dao
}

type Reader interface {
	First(dest interface{}, conds ...interface{}) Dao
}

type Dao interface {
	Creator
	Reader
	Error() error
	AutoMigrate(dst ...interface{}) error
}
