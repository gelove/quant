package orm

import (
	"quant/internal/app/entity"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// github.com/mattn/go-sqlite3
	_db, err := gorm.Open(sqlite.Open("quant.db"), &gorm.Config{})
	if err != nil {
		panic(errors.WithStack(err))
	}
	DB = _db
	err = DB.AutoMigrate(
		new(entity.Order),
	)
	if err != nil {
		panic(errors.WithStack(err))
	}
}
