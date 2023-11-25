package database

import (
	"gorm.io/gorm"
)

type Product struct {
	ID   uint   `gorm:"primary key;autoIncrement" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

func MigrateProducts(db *gorm.DB) error {
	err := db.AutoMigrate(&Product{})

	return err
}
