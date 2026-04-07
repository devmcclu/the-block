package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath+"?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Vehicle{}, &DamageNote{}, &VehicleImage{}, &Bid{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
