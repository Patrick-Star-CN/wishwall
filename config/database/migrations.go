package database

import (
	"gorm.io/gorm"
	"wishwall/app/models"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Wish{},
		&models.Stu{})
}
