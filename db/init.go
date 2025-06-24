// Package db contains stuff related to the database that's used to store shortened URL origins
package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ShortURL represents an entry in the database that includes the origin URL and it's
// shortened version
type ShortURL struct {
	Origin    string `gorm:"not null;unique"`
	Short     string `gorm:"not null;unique;index"`
	ExpiresAt int64  `gorm:"not null"`
}

// Init initializes an SQLite database
func Init() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	err = db.AutoMigrate(ShortURL{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database structs: %w", err)
	}

	return db, nil
}
