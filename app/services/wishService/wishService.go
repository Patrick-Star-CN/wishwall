package wishService

import (
	"wishwall/app/models"
	"wishwall/config/database"
)

func CreateWish(wish models.Wish) error {
	result := database.DB.Create(&wish)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetWishUser(uid int) ([]models.Wish, error) {
	var wish []models.Wish

	result := database.DB.Where(
		&models.Wish{
			UID: uid,
		}).Find(&wish)
	if result.Error != nil {
		return nil, result.Error
	}
	return wish, nil
}

func GetWishAll() ([]models.Wish, error) {
	var wish []models.Wish

	result := database.DB.Find(&wish)
	if result.Error != nil {
		return nil, result.Error
	}
	return wish, nil
}

func UpdateWish(wish models.Wish) error {
	result := database.DB.Model(models.Wish{}).Where(
		&models.Wish{
			ID: wish.ID,
		}).Updates(&wish)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteWish(id int) error {
	result := database.DB.Delete(&models.Wish{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
