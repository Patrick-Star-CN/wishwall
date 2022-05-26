package userService

import (
	"wishwall/app/models"
	"wishwall/config/database"
)

func CreateUser(user models.User) error {
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser(user models.User) error {
	result := database.DB.Model(models.User{}).Where(
		&models.User{
			ID: user.ID,
		}).Updates(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUser(username string) (*models.User, error) {
	var user *models.User

	result := database.DB.Where(
		&models.User{
			Name: username,
		}).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func GetUserID(id int) (*models.User, error) {
	var user *models.User

	result := database.DB.Where(
		&models.User{
			ID: id,
		}).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
