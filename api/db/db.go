package db

import (
	"github.com/jinzhu/gorm"
)

// DB returns gorm instance or error
func DB(dbname string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
