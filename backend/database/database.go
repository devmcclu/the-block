package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Vehicle{}, &DamageNote{}, &VehicleImage{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
