package kit

import (
	"gorm.io/gorm"
)

type Creator interface {
	Create(value interface{}) (tx *gorm.DB)
}

type Reader interface {
	Read(value interface{})
}
