package kit

type Creator interface {
	Create(value interface{}) DatabaseActions
}

type Reader interface {
	First(dest interface{}, conds ...interface{}) DatabaseActions
	Find(dest interface{}, conds ...interface{}) DatabaseActions
	Search(dest interface{}, query map[string]interface{}) DatabaseActions
}

type DatabaseActions interface {
	Creator
	Reader
	Error() error
	AutoMigrate(dst ...interface{}) error
}
